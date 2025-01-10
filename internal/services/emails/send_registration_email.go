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

//go:embed templates/registration_email.html
var templateRegistrationFile embed.FS

type templateRegistrationData struct {
	FrontURL                         string
	Token                            string
	FrontConfirmationRegistrationURL string
}

func (e *Emails) SendRegistrationEmail(
	inputCtx context.Context,
	email string,
	token string,
) error {
	tracer := otel.Tracer(tracing.GetServiceTraceName())
	ctx, span := tracer.Start(inputCtx, tracing.GetServiceTraceName())
	defer span.End()

	t, err := template.ParseFS(templateRegistrationFile, "templates/registration_email.html")
	if err != nil {
		tracing.HandleError(span, err)
		return errors.WithStack(err)
	}

	var message bytes.Buffer

	err = t.Execute(&message, templateRegistrationData{
		FrontURL:                         e.config.FrontURL,
		FrontConfirmationRegistrationURL: e.config.FrontConfirmationRegistrationURL,
		Token:                            token,
	})

	if err != nil {
		tracing.HandleError(span, err)
		return errors.WithStack(err)
	}

	err = e.repositories.ycp.SendEmail(
		ctx,
		email,
		fmt.Sprintf("Регистрация на %s", e.config.FrontURL),
		message.String(),
	)

	if err != nil {
		tracing.HandleError(span, err)
		return err
	}

	return nil
}
