package oauth

import (
	"github.com/upikoth/aireader-go/internal/config"
	"github.com/upikoth/aireader-go/internal/pkg/logger"
	"github.com/upikoth/aireader-go/internal/repositories/http/oauthmailru"
	"github.com/upikoth/aireader-go/internal/repositories/http/oauthyandex"
	"github.com/upikoth/aireader-go/internal/services/sessions"
	"github.com/upikoth/aireader-go/internal/services/users"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/mailru"
	"golang.org/x/oauth2/vk"
	"golang.org/x/oauth2/yandex"
)

type oauthRepositories struct {
	oauthMailRu *oauthmailru.OauthMailRu
	oauthYandex *oauthyandex.OauthYandex
}

type oauthServices struct {
	users    *users.Users
	sessions *sessions.Sessions
}

type Oauth struct {
	logger       logger.Logger
	vkConfig     *oauth2.Config
	mailConfig   *oauth2.Config
	yandexConfig *oauth2.Config
	repositories *oauthRepositories
	services     *oauthServices
}

func New(
	log logger.Logger,
	cfg config.Oauth,
	oauthMailRuRepo *oauthmailru.OauthMailRu,
	oauthYandexRepo *oauthyandex.OauthYandex,
	usersSrv *users.Users,
	sessionsSrv *sessions.Sessions,
) *Oauth {
	return &Oauth{
		logger: log,
		vkConfig: &oauth2.Config{
			ClientID:     cfg.VkClientID,
			ClientSecret: cfg.VkClientSecret,
			Endpoint:     vk.Endpoint,
			RedirectURL:  cfg.VkRedirectURL,
			Scopes:       []string{"email"},
		},
		mailConfig: &oauth2.Config{
			ClientID:     cfg.MailClientID,
			ClientSecret: cfg.MailClientSecret,
			Endpoint:     mailru.Endpoint,
			RedirectURL:  cfg.MailRedirectURL,
			Scopes:       []string{"userinfo"},
		},
		yandexConfig: &oauth2.Config{
			ClientID:     cfg.YandexClientID,
			ClientSecret: cfg.YandexClientSecret,
			Endpoint:     yandex.Endpoint,
			RedirectURL:  cfg.YandexRedirectURL,
		},
		repositories: &oauthRepositories{
			oauthMailRu: oauthMailRuRepo,
			oauthYandex: oauthYandexRepo,
		},
		services: &oauthServices{
			users:    usersSrv,
			sessions: sessionsSrv,
		},
	}
}
