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
func (o *Oauth) HandleMailRuRedirect(
	inputCtx context.Context,
	code string,
) (*models.SessionWithUserRole, error) {
	tracer := otel.Tracer(tracing.GetServiceTraceName())
	ctx, span := tracer.Start(inputCtx, tracing.GetServiceTraceName())
	defer span.End()

	token, err := o.mailConfig.Exchange(ctx, code)

	if err != nil {
		tracing.HandleError(span, err)
		return nil, err
	}

	userInfoMailRu, err := o.repositories.oauthMailRu.GetUserInfo(ctx, token.AccessToken)

	if err != nil {
		tracing.HandleError(span, err)
		return nil, err
	}

	span.SetAttributes(attribute.String("email", userInfoMailRu.Email))
	span.SetAttributes(attribute.String("userMailRuID", userInfoMailRu.ID))

	var user *models.User
	var errByMailRuID error
	var errByEmail error
	var errCreateUser error

	user, errByMailRuID = o.services.users.GetByMailRuID(ctx, userInfoMailRu.ID)

	// Если ошибка отличается от того что пользователь не найден.
	if errByMailRuID != nil && !errors.Is(errByMailRuID, constants.ErrUserNotFound) {
		tracing.HandleError(span, errByMailRuID)
		return nil, errByMailRuID
	}

	// Если пользователь не найден по mailru id.
	if errByMailRuID != nil && errors.Is(errByMailRuID, constants.ErrUserNotFound) {
		user, errByEmail = o.services.users.GetByEmail(ctx, userInfoMailRu.Email)

		// Если ошибка отличается от того что пользователь не найден.
		if errByEmail != nil && !errors.Is(errByEmail, constants.ErrUserNotFound) {
			tracing.HandleError(span, errByEmail)
			return nil, errByEmail
		}
	}

	// Если пользователь не найден даже по email, создаем нового.
	if errByEmail != nil && errors.Is(errByEmail, constants.ErrUserNotFound) {
		user, errCreateUser = o.services.users.CreateByEmailMailRuID(ctx, userInfoMailRu.Email, userInfoMailRu.ID)

		if errCreateUser != nil {
			tracing.HandleError(span, errCreateUser)
			return nil, errCreateUser
		}
	}

	// Если пользоваль найден по mailru id, обновляем почту если нужно, создаем сессию
	if errByMailRuID == nil {
		// Обновляем почту, если она отличается от той, что в токене.
		if user.Email != userInfoMailRu.Email {
			user.Email = userInfoMailRu.Email
			_, _ = o.services.users.UpdateUser(ctx, user)
		}
	}

	// Если пользоваль найден по email, обновляем mailru id если нужно, создаем сессию.
	if errByMailRuID != nil && errByEmail == nil {
		// Обновляем mailru id, если отличается от текущего.
		if user.MailRuID != userInfoMailRu.ID {
			user.MailRuID = userInfoMailRu.ID
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
