package ydbmodels

import "github.com/upikoth/aireader-go/internal/models"

type Session struct {
	ID       string
	Token    string
	UserID   string
	UserRole string
}

func NewYDBSessionModel(session *models.Session) *Session {
	return &Session{
		ID:     string(session.ID),
		Token:  session.Token,
		UserID: string(session.UserID),
	}
}

func (s *Session) FromYDBModel() *models.SessionWithUserRole {
	return &models.SessionWithUserRole{
		Session: models.Session{
			ID:     models.SessionID(s.ID),
			Token:  s.Token,
			UserID: models.UserID(s.UserID),
		},
		UserRole: models.UserRole(s.UserRole),
	}
}
