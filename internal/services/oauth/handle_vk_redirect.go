package oauth

import (
	"context"
	"fmt"

	"github.com/pkg/errors"
	"github.com/upikoth/aireader-go/internal/constants"
	"github.com/upikoth/aireader-go/internal/models"
	"github.com/upikoth/aireader-go/internal/pkg/tracing"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"golang.org/x/oauth2"
)

// nolint:funlen,gocognit,nolintlint //пока решил заигнорить.
func (o *Oauth) HandleVkRedirect(
	inputCtx context.Context,
	code string,
) (*models.SessionWithUserRole, error) {
	tracer := otel.Tracer(tracing.GetServiceTraceName())
	ctx, span := tracer.Start(inputCtx, tracing.GetServiceTraceName())
	defer span.End()

	token, err := o.vkConfig.Exchange(ctx, code)

	if err != nil {
		tracing.HandleError(span, err)
		return nil, err
	}

	email, err := getEmailFromVk(token)

	if err != nil {
		tracing.HandleError(span, err)
		return nil, err
	}

	userVkID, err := getUserVkID(token)

	if err != nil {
		tracing.HandleError(span, err)
		return nil, err
	}

	span.SetAttributes(attribute.String("email", email))
	span.SetAttributes(attribute.String("userVkID", userVkID))

	var user *models.User
	var errByVkID error
	var errByEmail error
	var errCreateUser error

	user, errByVkID = o.services.users.GetByVkID(ctx, userVkID)

	// Если ошибка отличается от того что пользователь не найден.
	if errByVkID != nil && !errors.Is(errByVkID, constants.ErrUserNotFound) {
		tracing.HandleError(span, errByVkID)
		return nil, errByVkID
	}

	// Если пользователь не найден по vk id.
	if errByVkID != nil && errors.Is(errByVkID, constants.ErrUserNotFound) {
		user, errByEmail = o.services.users.GetByEmail(ctx, email)

		// Если ошибка отличается от того что пользователь не найден.
		if errByEmail != nil && !errors.Is(errByEmail, constants.ErrUserNotFound) {
			tracing.HandleError(span, errByEmail)
			return nil, errByEmail
		}
	}

	// Если пользователь не найден даже по email, создаем нового.
	if errByEmail != nil && errors.Is(errByEmail, constants.ErrUserNotFound) {
		user, errCreateUser = o.services.users.CreateByEmailVkID(ctx, email, userVkID)

		if errCreateUser != nil {
			tracing.HandleError(span, errCreateUser)
			return nil, errCreateUser
		}
	}

	// Если пользоваль найден по vk id, обновляем почту если нужно, создаем сессию
	if errByVkID == nil {
		// Обновляем почту, если она отличается от той, что в токене.
		if user.Email != email {
			user.Email = email
			_, _ = o.services.users.UpdateUser(ctx, user)
		}
	}

	// Если пользоваль найден по email, обновляем vk id если нужно, создаем сессию.
	if errByVkID != nil && errByEmail == nil {
		// Обновляем vk id, если отличается от текущего.
		if user.VkID != userVkID {
			user.VkID = userVkID
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

func getEmailFromVk(token *oauth2.Token) (string, error) {
	email := token.Extra("email")

	parsedEmail, emailOk := email.(string)

	if !emailOk || parsedEmail == "" {
		return "", errors.New("Invalid email")
	}

	return parsedEmail, nil
}

func getUserVkID(token *oauth2.Token) (string, error) {
	vkUserID := token.Extra("user_id")

	userIDInt, userIDOk := vkUserID.(float64)

	if !userIDOk || userIDInt == 0 {
		return "", errors.New("Invalid user ID")
	}

	parsedUserID := fmt.Sprintf("%.0f", userIDInt)

	return parsedUserID, nil
}
