package handler

import (
	"context"
	"net/http"

	"github.com/pkg/errors"
	"github.com/upikoth/aireader-go/internal/constants"
	app "github.com/upikoth/aireader-go/internal/generated/app"
	"github.com/upikoth/aireader-go/internal/models"
	"github.com/upikoth/aireader-go/internal/pkg/tracing"
	"go.opentelemetry.io/otel"
)

func (h *Handler) V1GetUsers(
	inputCtx context.Context,
	params app.V1GetUsersParams,
) (*app.V1UsersGetUsersResponse, error) {
	tracer := otel.Tracer(tracing.GetHandlerTraceName())
	ctx, span := tracer.Start(inputCtx, tracing.GetHandlerTraceName())
	defer span.End()

	session, err := h.services.Sessions.CheckToken(ctx, params.AuthorizationToken)

	if errors.Is(err, constants.ErrSessionNotFound) {
		return nil, &models.Error{
			Code:        models.ErrCodeUserUnauthorized,
			Description: "User session is invalid",
			StatusCode:  http.StatusUnauthorized,
		}
	}

	if err != nil {
		tracing.HandleError(span, err)
		return nil, &models.Error{
			Code:        models.ErrCodeInterval,
			Description: err.Error(),
		}
	}

	if !session.UserRole.CheckAccessToAction(models.UserActionGetAnyUserInfo) {
		return nil, &models.Error{
			Code:        models.ErrCodeUsersGetListForbidden,
			Description: "Insufficient rights",
			StatusCode:  http.StatusForbidden,
		}
	}

	usersGetListParams := &models.UsersGetListParams{
		Limit:  params.Limit.Value,
		Offset: params.Offset.Value,
	}

	userList, err := h.services.Users.GetList(ctx, usersGetListParams)

	if err != nil {
		tracing.HandleError(span, err)
		return nil, &models.Error{
			Code:        models.ErrCodeInterval,
			Description: err.Error(),
		}
	}

	var usersResult []app.User
	for _, user := range userList.Users {
		usersResult = append(usersResult, app.User{
			ID:    string(user.ID),
			Email: user.Email,
			Role:  app.UserRole(user.Role),
		})
	}

	return &app.V1UsersGetUsersResponse{
		Success: true,
		Data: app.V1UsersGetUsersResponseData{
			Users:  usersResult,
			Limit:  usersGetListParams.Limit,
			Offset: usersGetListParams.Offset,
			Total:  userList.Total,
		},
	}, nil
}

func (h *Handler) V1GetCurrentUser(
	inputCtx context.Context,
	params app.V1GetCurrentUserParams,
) (*app.V1UsersGetUserResponse, error) {
	tracer := otel.Tracer(tracing.GetHandlerTraceName())
	ctx, span := tracer.Start(inputCtx, tracing.GetHandlerTraceName())
	defer span.End()

	session, err := h.services.Sessions.CheckToken(ctx, params.AuthorizationToken)

	if errors.Is(err, constants.ErrSessionNotFound) {
		return nil, &models.Error{
			Code:        models.ErrCodeUserUnauthorized,
			Description: "User session is invalid",
			StatusCode:  http.StatusUnauthorized,
		}
	}

	if err != nil {
		tracing.HandleError(span, err)
		return nil, &models.Error{
			Code:        models.ErrCodeInterval,
			Description: err.Error(),
		}
	}

	user, err := h.services.Users.GetByID(ctx, session.UserID)

	if err != nil {
		tracing.HandleError(span, err)
		return nil, &models.Error{
			Code:        models.ErrCodeInterval,
			Description: err.Error(),
		}
	}

	return &app.V1UsersGetUserResponse{
		Success: true,
		Data: app.V1UsersGetUserResponseData{
			User: app.User{
				ID:    string(user.ID),
				Email: user.Email,
				Role:  app.UserRole(user.Role),
			},
		},
	}, nil
}
