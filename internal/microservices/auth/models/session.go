package models

import (
	"context"
	"crypto/md5"
	"dripapp/internal/dripapp/models"
	"fmt"
	"net/http"
	"time"
)

//easyjson:json
type Session struct {
	Cookie string
	UserID uint64
}

func CreateSessionCookie(user models.LoginUser) http.Cookie {
	expiration := time.Now().Add(10 * time.Hour)

	data := user.Password + time.Now().String()
	md5CookieValue := fmt.Sprintf("%x", md5.Sum([]byte(data)))

	cookie := http.Cookie{
		Name:     "sessionId",
		Value:    md5CookieValue,
		Expires:  expiration,
		Secure:   true,
		HttpOnly: true,
		SameSite: http.SameSiteNoneMode,
		Path:     "/api/v1",
	}

	return cookie
}

func NewSession(id uint64, email, cookie string) *Session {
	return &Session{UserID: id, Cookie: cookie}
}

type SessionRepository interface {
	GetSessionByCookie(sessionCookie string) (Session, error)
	NewSessionCookie(sessionCookie string, userId uint64) error
	DeleteSessionCookie(sessionCookie string) error
}

type SessionUsecase interface {
	AddSession(c context.Context, session Session) error
	DeleteSession(c context.Context) error
}
