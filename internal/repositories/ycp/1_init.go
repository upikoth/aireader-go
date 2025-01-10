package ycp

import (
	"github.com/pkg/errors"
	"github.com/upikoth/aireader-go/internal/config"
	"github.com/upikoth/aireader-go/internal/pkg/logger"
	smtpclient "github.com/upikoth/aireader-go/internal/pkg/smtp-client"
)

type Ycp struct {
	logger logger.Logger
	config *config.Ycp
	client *smtpclient.SMTPClient
}

func New(
	logger logger.Logger,
	config *config.Ycp,
) (*Ycp, error) {
	client, err := smtpclient.New(
		config.Host,
		config.Port,
		config.Username,
		config.Password,
	)

	if err != nil {
		return nil, errors.WithStack(err)
	}

	return &Ycp{
		logger: logger,
		config: config,
		client: client,
	}, nil
}
