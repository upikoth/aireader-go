package services

import (
	"github.com/upikoth/aireader-go/internal/config"
	"github.com/upikoth/aireader-go/internal/pkg/logger"
	"github.com/upikoth/aireader-go/internal/repositories"
	"github.com/upikoth/aireader-go/internal/services/emails"
	"github.com/upikoth/aireader-go/internal/services/oauth"
	passwordrecoveryrequests "github.com/upikoth/aireader-go/internal/services/password-recovery-requests"
	"github.com/upikoth/aireader-go/internal/services/registrations"
	"github.com/upikoth/aireader-go/internal/services/sessions"
	"github.com/upikoth/aireader-go/internal/services/users"
)

type Services struct {
	Registrations            *registrations.Registrations
	Sessions                 *sessions.Sessions
	PasswordRecoveryRequests *passwordrecoveryrequests.PasswordRecoveryRequests
	Users                    *users.Users
	Oauth                    *oauth.Oauth
	Emails                   *emails.Emails
}

func New(
	log logger.Logger,
	cfg *config.Config,
	repo *repositories.Repository,
) (*Services, error) {
	srvs := &Services{}

	srvs.Emails = emails.New(
		log,
		&cfg.Services.Emails,
		repo.YCP,
	)

	srvs.Users = users.New(
		log,
		repo.YDB.DB,
		repo.YDB.Users,
	)

	srvs.Sessions = sessions.New(
		log,
		repo.YDB.DB,
		repo.YDB.Sessions,
		srvs.Users,
	)

	srvs.Registrations = registrations.New(
		log,
		repo.YDB.DB,
		repo.YDB.Registrations,
		srvs.Users,
		srvs.Sessions,
		srvs.Emails,
	)

	srvs.PasswordRecoveryRequests = passwordrecoveryrequests.New(
		log,
		repo.YDB.DB,
		repo.YDB.PasswordRecoveryRequests,
		srvs.Users,
		srvs.Sessions,
		srvs.Emails,
	)

	srvs.Oauth = oauth.New(
		log,
		cfg.Services.Oauth,
		repo.HTTP.OauthMailRu,
		repo.HTTP.OauthYandex,
		srvs.Users,
		srvs.Sessions,
	)

	return srvs, nil
}
