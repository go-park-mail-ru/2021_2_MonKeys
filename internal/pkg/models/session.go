package models

type Session struct {
	UserID uint64
	Cookie string
}

func NewSession(id uint64, email, cookie string) *Session {
	return &Session{UserID: id, Cookie: cookie}
}

type SessionRepository interface {
	GetUserIDByCookie(sessionCookie string) (uint64, error)
	NewSessionCookie(sessionCookie string, userId uint64) error
	DeleteSessionCookie(sessionCookie string) error
	IsSessionByUserID(userID uint64) bool
	DropCookies()
}
