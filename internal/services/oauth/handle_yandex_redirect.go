package oauth

import (
	"context"

	"github.com/pkg/errors"
	"github.com/upikoth/aireader-go/internal/constants"
	"github.com/upikoth/aireader-go/internal/models"
	"github.com/upikoth/aireader-go/internal/pkg/tracing"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
)

// nolint:funlen,gocognit,nolintlint //пока решил заигнорить.
func (o *Oauth) HandleYandexRedirect(
	inputCtx context.Context,
	code string,
) (*models.SessionWithUserRole, error) {
	tracer := otel.Tracer(tracing.GetServiceTraceName())
	ctx, span := tracer.Start(inputCtx, tracing.GetServiceTraceName())
	defer span.End()

	token, err := o.yandexConfig.Exchange(ctx, code)

	if err != nil {
		tracing.HandleError(span, err)
		return nil, err
	}

	userInfoYandex, err := o.repositories.oauthYandex.GetUserInfo(ctx, token.AccessToken)

	if err != nil {
		tracing.HandleError(span, err)
		return nil, err
	}

	span.SetAttributes(attribute.String("email", userInfoYandex.Email))
	span.SetAttributes(attribute.String("userYandexID", userInfoYandex.ID))

	var user *models.User
	var errByYandexID error
	var errByEmail error
	var errCreateUser error

	user, errByYandexID = o.services.users.GetByYandexID(ctx, userInfoYandex.ID)

	// Если ошибка отличается от того что пользователь не найден.
	if errByYandexID != nil && !errors.Is(errByYandexID, constants.ErrUserNotFound) {
		tracing.HandleError(span, errByYandexID)
		return nil, errByYandexID
	}

	// Если пользователь не найден по yandex id.
	if errByYandexID != nil && errors.Is(errByYandexID, constants.ErrUserNotFound) {
		user, errByEmail = o.services.users.GetByEmail(ctx, userInfoYandex.Email)

		// Если ошибка отличается от того что пользователь не найден.
		if errByEmail != nil && !errors.Is(errByEmail, constants.ErrUserNotFound) {
			tracing.HandleError(span, errByEmail)
			return nil, errByEmail
		}
	}

	// Если пользователь не найден даже по email, создаем нового.
	if errByEmail != nil && errors.Is(errByEmail, constants.ErrUserNotFound) {
		user, errCreateUser = o.services.users.CreateByEmailYandexID(ctx, userInfoYandex.Email, userInfoYandex.ID)

		if errCreateUser != nil {
			tracing.HandleError(span, errCreateUser)
			return nil, errCreateUser
		}
	}

	// Если пользоваль найден по yandex id, обновляем почту если нужно, создаем сессию
	if errByYandexID == nil {
		// Обновляем почту, если она отличается от той, что в токене.
		if user.Email != userInfoYandex.Email {
			user.Email = userInfoYandex.Email
			_, _ = o.services.users.UpdateUser(ctx, user)
		}
	}

	// Если пользоваль найден по email, обновляем yandex id если нужно, создаем сессию.
	if errByYandexID != nil && errByEmail == nil {
		// Обновляем yandex id, если отличается от текущего.
		if user.YandexID != userInfoYandex.ID {
			user.YandexID = userInfoYandex.ID
			_, _ = o.services.users.UpdateUser(ctx, user)
		}
	}

	session, errS := o.services.sessions.CreateByUserID(ctx, user.ID)

	if errS != nil {
		tracing.HandleError(span, errS)
		return nil, errS
	}

	return session, nil
}
