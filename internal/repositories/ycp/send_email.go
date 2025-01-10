package ycp

import (
	"context"

	"github.com/getsentry/sentry-go"
	"github.com/pkg/errors"
	"github.com/upikoth/aireader-go/internal/pkg/tracing"
	"go.opentelemetry.io/otel"

	"net/mail"
)

func (y *Ycp) SendEmail(
	inputCtx context.Context,
	toEmail string,
	title string,
	body string,
) (err error) {
	tracer := otel.Tracer(tracing.GetRepositoryTraceName())
	_, span := tracer.Start(inputCtx, tracing.GetRepositoryTraceName())
	defer span.End()

	defer func() {
		if err != nil {
			span.RecordError(err)
			sentry.CaptureException(err)
		}
	}()

	err = y.client.Connect()

	if err != nil {
		return errors.WithStack(err)
	}

	defer func() {
		disconnectError := y.client.Disconnect()

		if err != nil {
			y.logger.Error(disconnectError.Error())
		}
	}()

	from := mail.Address{
		Name:    y.config.FromName,
		Address: y.config.FromAddress,
	}
	to := mail.Address{
		Name:    "",
		Address: toEmail,
	}

	message := y.client.CreateMessage(from, to, title, body)
	err = y.client.SendEmail(from, to, message)

	return errors.WithStack(err)
}
