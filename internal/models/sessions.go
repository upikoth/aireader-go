package models

type SessionID string

type Session struct {
	ID     SessionID
	Token  string
	UserID UserID
}

type SessionWithUserRole struct {
	Session
	UserRole UserRole
}
