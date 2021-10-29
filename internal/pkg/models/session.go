package models

type Session struct {
	Cookie string
	UserID uint64
}

func NewSession(id uint64, email, cookie string) *Session {
	return &Session{UserID: id, Cookie: cookie}
}

type SessionRepository interface {
	GetSessionByCookie(sessionCookie string) (Session, error)
	NewSessionCookie(sessionCookie string, userId uint64) error
	DeleteSessionCookie(sessionCookie string) error
	IsSessionByCookie(sessionCookie string) bool
	DropCookies()
}
