package usecase

import (
	"context"
	"dripapp/configs"
	_sessionModels "dripapp/internal/microservices/auth/models"
	"errors"
	"time"
)

type sessionUsecase struct {
	Session        _sessionModels.SessionRepository
	contextTimeout time.Duration
}

func NewSessionUsecase(sess _sessionModels.SessionRepository, timeout time.Duration) _sessionModels.SessionUsecase {
	return &sessionUsecase{
		Session:        sess,
		contextTimeout: timeout,
	}
}

func (s *sessionUsecase) AddSession(c context.Context, session _sessionModels.Session) error {
	err := s.Session.NewSessionCookie(session.Cookie, session.UserID)
	if err != nil {
		return err
	}
	return nil
}

func (s *sessionUsecase) DeleteSession(c context.Context) error {
	ctx, cancel := context.WithTimeout(c, s.contextTimeout)
	defer cancel()

	ctxSession := ctx.Value(configs.ContextUserID)
	if ctxSession == nil {
		return errors.New("context nil error")
	}
	currentSession, ok := ctxSession.(_sessionModels.Session)
	if !ok {
		return errors.New("convert to model session error")
	}
	err := s.Session.DeleteSessionCookie(currentSession.Cookie)
	if err != nil {
		return err
	}
	return nil
}
