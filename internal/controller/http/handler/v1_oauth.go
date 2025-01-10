package handler

import (
	"context"
	"fmt"
	"net/http"

	app "github.com/upikoth/aireader-go/internal/generated/app"
	"github.com/upikoth/aireader-go/internal/models"
	"github.com/upikoth/aireader-go/internal/pkg/tracing"
	"go.opentelemetry.io/otel"
)

func (h *Handler) V1AuthorizeUsingOauth(
	inputCtx context.Context,
	req *app.V1AuthorizeUsingOauthRequestBody,
) (*app.V1AuthorizeUsingOauthResponse, error) {
	tracer := otel.Tracer(tracing.GetHandlerTraceName())
	ctx, span := tracer.Start(inputCtx, tracing.GetHandlerTraceName())
	defer span.End()

	var oauthSource models.OauthSource

	switch req.OauthSource {
	case app.V1AuthorizeUsingOauthRequestBodyOauthSourceVk:
		oauthSource = models.OauthSourceVk
	case app.V1AuthorizeUsingOauthRequestBodyOauthSourceMail:
		oauthSource = models.OauthSourceMail
	case app.V1AuthorizeUsingOauthRequestBodyOauthSourceYandex:
		oauthSource = models.OauthSourceYandex
	default:
		return nil, &models.Error{
			Code:        models.ErrCodeOauthSourceNotExist,
			Description: "Incorrect email or password",
			StatusCode:  http.StatusBadRequest,
		}
	}

	url, err := h.services.Oauth.GetAuthorizeURL(ctx, oauthSource)

	if err != nil {
		return nil, err
	}

	return &app.V1AuthorizeUsingOauthResponse{
		Success: true,
		Data: app.V1AuthorizeUsingOauthResponseData{
			URL: url,
		},
	}, nil
}

func (h *Handler) V1AuthorizeUsingOauthHandleVkRedirect(
	inputCtx context.Context,
	params app.V1AuthorizeUsingOauthHandleVkRedirectParams,
) (*app.V1AuthorizeUsingOauthHandleVkRedirectFound, error) {
	tracer := otel.Tracer(tracing.GetHandlerTraceName())
	ctx, span := tracer.Start(inputCtx, tracing.GetHandlerTraceName())
	defer span.End()

	session, err := h.services.Oauth.HandleVkRedirect(ctx, params.Code)

	if err != nil {
		return nil, &models.Error{
			Code:        models.ErrCodeInterval,
			Description: err.Error(),
		}
	}

	return &app.V1AuthorizeUsingOauthHandleVkRedirectFound{
		Location: app.NewOptString(
			fmt.Sprintf(
				"%s?id=%s&token=%s&userRole=%s",
				h.cfg.FrontHandleAuthPageURL,
				string(session.ID),
				session.Token,
				string(session.UserRole),
			),
		),
	}, nil
}

func (h *Handler) V1AuthorizeUsingOauthHandleMailRedirect(
	inputCtx context.Context,
	params app.V1AuthorizeUsingOauthHandleMailRedirectParams,
) (*app.V1AuthorizeUsingOauthHandleMailRedirectFound, error) {
	tracer := otel.Tracer(tracing.GetHandlerTraceName())
	ctx, span := tracer.Start(inputCtx, tracing.GetHandlerTraceName())
	defer span.End()

	session, err := h.services.Oauth.HandleMailRuRedirect(ctx, params.Code)

	if err != nil {
		return nil, &models.Error{
			Code:        models.ErrCodeInterval,
			Description: err.Error(),
		}
	}

	return &app.V1AuthorizeUsingOauthHandleMailRedirectFound{
		Location: app.NewOptString(
			fmt.Sprintf(
				"%s?id=%s&token=%s&userRole=%s",
				h.cfg.FrontHandleAuthPageURL,
				string(session.ID),
				session.Token,
				string(session.UserRole),
			),
		),
	}, nil
}

func (h *Handler) V1AuthorizeUsingOauthHandleYandexRedirect(
	inputCtx context.Context,
	params app.V1AuthorizeUsingOauthHandleYandexRedirectParams,
) (*app.V1AuthorizeUsingOauthHandleYandexRedirectFound, error) {
	tracer := otel.Tracer(tracing.GetHandlerTraceName())
	ctx, span := tracer.Start(inputCtx, tracing.GetHandlerTraceName())
	defer span.End()

	session, err := h.services.Oauth.HandleYandexRedirect(ctx, params.Code)

	if err != nil {
		return nil, &models.Error{
			Code:        models.ErrCodeInterval,
			Description: err.Error(),
		}
	}

	return &app.V1AuthorizeUsingOauthHandleYandexRedirectFound{
		Location: app.NewOptString(
			fmt.Sprintf(
				"%s?id=%s&token=%s&userRole=%s",
				h.cfg.FrontHandleAuthPageURL,
				string(session.ID),
				session.Token,
				string(session.UserRole),
			),
		),
	}, nil
}
