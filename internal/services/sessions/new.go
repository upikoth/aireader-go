package sessions

import (
	"github.com/google/uuid"
	"github.com/upikoth/aireader-go/internal/models"
)

type Option func(session *models.Session)

func newSession(
	userID models.UserID,
	options ...Option,
) *models.Session {
	session := &models.Session{
		ID:     models.SessionID(uuid.New().String()),
		Token:  uuid.New().String(),
		UserID: userID,
	}

	for _, option := range options {
		option(session)
	}

	return session
}
