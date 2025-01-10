package emails

import (
	"bytes"
	"context"
	"embed"
	"fmt"
	"html/template"

	"github.com/pkg/errors"
	"github.com/upikoth/aireader-go/internal/pkg/tracing"
	"go.opentelemetry.io/otel"
)

//go:embed templates/password_recovery_request_email.html
var templateFile embed.FS

type templateData struct {
	FrontURL                                    string
	Token                                       string
	FrontConfirmationPasswordRecoveryRequestURL string
}

func (e *Emails) SendPasswordRecoveryRequestEmail(
	inputCtx context.Context,
	email string,
	token string,
) error {
	tracer := otel.Tracer(tracing.GetServiceTraceName())
	ctx, span := tracer.Start(inputCtx, tracing.GetServiceTraceName())
	defer span.End()

	t, err := template.ParseFS(templateFile, "templates/password_recovery_request_email.html")
	if err != nil {
		tracing.HandleError(span, err)
		return errors.WithStack(err)
	}

	var message bytes.Buffer

	err = t.Execute(&message, templateData{
		FrontURL: e.config.FrontURL,
		FrontConfirmationPasswordRecoveryRequestURL: e.config.FrontConfirmationPasswordRecoveryRequestURL,
		Token: token,
	})

	if err != nil {
		tracing.HandleError(span, err)
		return errors.WithStack(err)
	}

	err = e.repositories.ycp.SendEmail(
		ctx,
		email,
		fmt.Sprintf("Восстановление пароля на %s", e.config.FrontURL),
		message.String(),
	)

	if err != nil {
		tracing.HandleError(span, err)
		return err
	}

	return nil
}
