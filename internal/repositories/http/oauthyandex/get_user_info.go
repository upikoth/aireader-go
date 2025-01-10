package oauthyandex

import (
	"context"

	"github.com/pkg/errors"
	oauthyandex "github.com/upikoth/aireader-go/internal/generated/oauthyandex"
	"github.com/upikoth/aireader-go/internal/models"
	"github.com/upikoth/aireader-go/internal/pkg/tracing"
	"go.opentelemetry.io/otel"
)

func (o *OauthYandex) GetUserInfo(
	inputCtx context.Context,
	accessToken string,
) (*models.OauthUserInfo, error) {
	tracer := otel.Tracer(tracing.GetRepositoryHTTPTraceName())
	ctx, span := tracer.Start(inputCtx, tracing.GetRepositoryHTTPTraceName())
	defer span.End()

	res, err := o.client.UserInfo(ctx, oauthyandex.UserInfoParams{
		OAuthToken: accessToken,
		Format:     oauthyandex.FormatJSON,
	})

	if err != nil {
		return nil, errors.WithStack(err)
	}

	return &models.OauthUserInfo{
		ID:    res.ID,
		Email: res.DefaultEmail,
	}, nil
}
