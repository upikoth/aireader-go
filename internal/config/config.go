package config

import (
	"github.com/kelseyhightower/envconfig"
	"github.com/pkg/errors"
)

type Config struct {
	Environment  string `envconfig:"ENVIRONMENT" required:"true"`
	Controller   Controller
	Services     Services
	Repositories Repositories
}

type Controller struct {
	HTTP ControllerHTTP
}

type ControllerHTTP struct {
	Port                   string `envconfig:"PORT" required:"true"`
	SentryDsn              string `envconfig:"SENTRY_DSN" required:"true"`
	Environment            string `envconfig:"ENVIRONMENT" required:"true"`
	FrontHandleAuthPageURL string `envconfig:"FRONT_HANDLE_AUTH_PAGE_URL" required:"true"`
}

type Services struct {
	Emails Emails
	Oauth  Oauth
}

type Emails struct {
	FrontURL                         string `envconfig:"FRONT_URL" required:"true"`
	FrontConfirmationRegistrationURL string `envconfig:"FRONT_CONFIRMATION_REGISTRATION_URL" required:"true"`
	//nolint:lll
	FrontConfirmationPasswordRecoveryRequestURL string `envconfig:"FRONT_CONFIRMATION_PASSWORD_RECOVERY_REQUEST_URL" required:"true"`
}

type Oauth struct {
	VkClientID     string `envconfig:"OAUTH_VK_CLIENT_ID" required:"true"`
	VkClientSecret string `envconfig:"OAUTH_VK_CLIENT_SECRET" required:"true"`
	VkRedirectURL  string `envconfig:"OAUTH_VK_REDIRECT_URL" required:"true"`

	MailClientID     string `envconfig:"OAUTH_MAIL_CLIENT_ID" required:"true"`
	MailClientSecret string `envconfig:"OAUTH_MAIL_CLIENT_SECRET" required:"true"`
	MailRedirectURL  string `envconfig:"OAUTH_MAIL_REDIRECT_URL" required:"true"`

	YandexClientID     string `envconfig:"OAUTH_YANDEX_CLIENT_ID" required:"true"`
	YandexClientSecret string `envconfig:"OAUTH_YANDEX_CLIENT_SECRET" required:"true"`
	YandexRedirectURL  string `envconfig:"OAUTH_YANDEX_REDIRECT_URL" required:"true"`
}

type Repositories struct {
	Ycp  Ycp
	Ydb  Ydb
	HTTP HTTP
}

type Ycp struct {
	Host        string `envconfig:"YCP_HOST" required:"true"`
	Port        string `envconfig:"YCP_PORT" required:"true"`
	FromName    string `envconfig:"YCP_FROM_NAME" required:"true"`
	FromAddress string `envconfig:"YCP_FROM_ADDRESS" required:"true"`
	Username    string `envconfig:"YCP_USERNAME" required:"true"`
	Password    string `envconfig:"YCP_PASSWORD" required:"true"`
}

type Ydb struct {
	Dsn                 string `envconfig:"YDB_DSN" required:"true"`
	AuthFileDirName     string `envconfig:"YDB_AUTH_FILE_DIR_NAME" required:"true"`
	AuthFileName        string `envconfig:"YDB_AUTH_FILE_NAME" required:"true"`
	YcSaJSONCredentials []byte `envconfig:"YC_SA_JSON_CREDENTIALS"`
	Environment         string `envconfig:"ENVIRONMENT" required:"true"`
}

type HTTP struct {
	OauthMailRu OauthMailRu
	OauthYandex OauthYandex
}

type OauthMailRu struct {
	APIURL string `envconfig:"OAUTH_MAIL_API_URL" required:"true"`
}

type OauthYandex struct {
	APIURL string `envconfig:"OAUTH_YANDEX_API_URL" required:"true"`
}

func New() (*Config, error) {
	config := &Config{}

	if err := envconfig.Process("", config); err != nil {
		return nil, errors.WithStack(err)
	}

	return config, nil
}
