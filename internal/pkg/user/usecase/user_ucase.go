package usecase

import (
	"dripapp/internal/pkg/models"
	"time"
)

type userUsecase struct {
	userRepo       models.UserRepository
	contextTimeout time.Duration
}
