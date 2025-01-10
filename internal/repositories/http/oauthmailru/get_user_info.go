package oauthmailru

import (
	"context"

	"github.com/pkg/errors"
	oauthmailru "github.com/upikoth/aireader-go/internal/generated/oauthmailru"
	"github.com/upikoth/aireader-go/internal/models"
	"github.com/upikoth/aireader-go/internal/pkg/tracing"
	"go.opentelemetry.io/otel"
)

func (o *OauthMailRu) GetUserInfo(
	inputCtx context.Context,
	accessToken string,
) (*models.OauthUserInfo, error) {
	tracer := otel.Tracer(tracing.GetRepositoryHTTPTraceName())
	ctx, span := tracer.Start(inputCtx, tracing.GetRepositoryHTTPTraceName())
	defer span.End()

	res, err := o.client.UserInfo(ctx, oauthmailru.UserInfoParams{
		AccessToken: accessToken,
	})

	if err != nil {
		return nil, errors.WithStack(err)
	}

	return &models.OauthUserInfo{
		ID:    res.ID,
		Email: res.Email,
	}, nil
}
