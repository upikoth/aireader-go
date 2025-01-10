package oauth

import (
	"context"
	"net/http"

	"github.com/google/uuid"
	"github.com/upikoth/aireader-go/internal/models"
	"github.com/upikoth/aireader-go/internal/pkg/tracing"
	"go.opentelemetry.io/otel"
)

func (o *Oauth) GetAuthorizeURL(
	inputCtx context.Context,
	oauthSource models.OauthSource,
) (string, error) {
	tracer := otel.Tracer(tracing.GetServiceTraceName())
	_, span := tracer.Start(inputCtx, tracing.GetServiceTraceName())
	defer span.End()

	switch oauthSource {
	case models.OauthSourceVk:
		return o.vkConfig.AuthCodeURL(uuid.New().String()), nil
	case models.OauthSourceMail:
		return o.mailConfig.AuthCodeURL(uuid.New().String()), nil
	case models.OauthSourceYandex:
		return o.yandexConfig.AuthCodeURL(uuid.New().String()), nil
	default:
		return "", &models.Error{
			Code:        models.ErrCodeOauthSourceNotExist,
			Description: "Incorrect oauth source",
			StatusCode:  http.StatusBadRequest,
		}
	}
}
