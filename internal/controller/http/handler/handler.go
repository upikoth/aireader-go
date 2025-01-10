package handler

import (
	"context"
	"errors"
	"net/http"

	"github.com/upikoth/aireader-go/internal/config"
	app "github.com/upikoth/aireader-go/internal/generated/app"
	"github.com/upikoth/aireader-go/internal/models"
	"github.com/upikoth/aireader-go/internal/pkg/logger"
	"github.com/upikoth/aireader-go/internal/services"
)

type Handler struct {
	logger   logger.Logger
	services *services.Services
	cfg      *config.ControllerHTTP
}

func New(
	log logger.Logger,
	srv *services.Services,
	cfg *config.ControllerHTTP,
) *Handler {
	return &Handler{
		logger:   log,
		services: srv,
		cfg:      cfg,
	}
}

func (h *Handler) NewError(_ context.Context, err error) *app.ErrorResponseStatusCode {
	modelErr := &models.Error{}

	code := models.ErrCodeValidationByOpenapi
	statusCode := http.StatusBadRequest

	if errors.As(err, &modelErr) {
		code = modelErr.Code
		statusCode = modelErr.GetStatusCode()
	}

	return &app.ErrorResponseStatusCode{
		StatusCode: statusCode,
		Response: app.ErrorResponse{
			Success: app.ErrorResponseSuccessFalse,
			Error: app.ErrorResponseError{
				Code:        string(code),
				Description: err.Error(),
			},
		},
	}
}
