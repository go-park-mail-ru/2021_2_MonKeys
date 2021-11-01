package usecase

import (
	"dripapp/internal/dripapp/models"

	_ "github.com/golang/mock/mockgen/model"
)

type SessionManager interface {
	GetSessionByCookie(string) (session models.Session, err error)
	NewSessionCookie(sessionCookie string, id uint64) error
	DeleteSessionCookie(sessionCookie string) error
	IsSessionByCookie(sessionCookie string) bool
}
