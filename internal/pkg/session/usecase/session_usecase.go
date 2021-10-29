package usecase

import (
	"context"
	"dripapp/internal/pkg/models"
	"errors"
	"time"
)

type sessionUsecase struct {
	Session        models.SessionRepository
	contextTimeout time.Duration
}

func NewSessionUsecase(sess models.SessionRepository, timeout time.Duration) models.SessionUsecase {
	return &sessionUsecase{
		Session:        sess,
		contextTimeout: timeout,
	}
}

func (s *sessionUsecase) AddSession(c context.Context, session models.Session) error {
	err := s.Session.NewSessionCookie(session.Cookie, session.UserID)
	if err != nil {
		return err
	}
	return nil
}

func (s *sessionUsecase) DeleteSession(c context.Context) error {
	ctx, cancel := context.WithTimeout(c, s.contextTimeout)
	defer cancel()

	ctxSession := ctx.Value("userID")
	if ctxSession == nil {
		return errors.New("context nil error")
	}
	currentSession, ok := ctxSession.(models.Session)
	if !ok {
		return errors.New("convert to model session error")
	}
	err := s.Session.DeleteSessionCookie(currentSession.Cookie)
	if err != nil {
		return err
	}
	return nil
}
