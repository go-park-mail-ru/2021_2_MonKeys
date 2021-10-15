package delivery

import (
	"dripapp/internal/pkg/models"
	"html/template"

	"dripapp/internal/pkg/session"

	"go.uber.org/zap"
)

type UserHandler struct {
	Tmpl      *template.Template
	Logger    *zap.SugaredLogger
	UserUCase models.UserUsecase
	Sessions  *session.SessionManager
}
