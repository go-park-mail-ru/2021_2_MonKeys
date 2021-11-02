package usecase

import (
	"context"
	"dripapp/configs"
	"dripapp/internal/dripapp/models"
	"dripapp/internal/pkg/hasher"
	"io"
	"net/http"
	"strconv"
	"time"
)

type userUsecase struct {
	UserRepo       models.UserRepository
	Session        models.SessionRepository
	File           models.FileRepository
	contextTimeout time.Duration
}

const (
	StatusOK                  = 200
	StatusBadRequest          = 400
	StatusNotFound            = 404
	StatusInternalServerError = 500
	StatusEmailAlreadyExists  = 1001
)

func NewUserUsecase(
	ur models.UserRepository,
	fileManager models.FileRepository,
	timeout time.Duration) models.UserUsecase {
	return &userUsecase{
		UserRepo:       ur,
		File:           fileManager,
		contextTimeout: timeout,
	}
}

func (h *userUsecase) CurrentUser(c context.Context) (models.User, models.HTTPError) {
	ctx, cancel := context.WithTimeout(c, h.contextTimeout)
	defer cancel()

	ctxSession := ctx.Value(configs.ForContext)
	if ctxSession == nil {
		return models.User{}, models.HTTPError{
			Code:    http.StatusNotFound,
			Message: models.ErrContextNilError,
		}
	}
	currentSession, ok := ctxSession.(models.Session)
	if !ok {
		return models.User{}, models.HTTPError{
			Code:    http.StatusNotFound,
			Message: models.ErrConvertToSession,
		}
	}

	currentUser, err := h.UserRepo.GetUserByID(c, currentSession.UserID)
	if err != nil {
		return models.User{}, models.HTTPError{
			Code:    http.StatusNotFound,
			Message: err.Error(),
		}
	}

	return currentUser, models.StatusOk200
}

func (h *userUsecase) EditProfile(c context.Context, newUserData models.User) (models.User, models.HTTPError) {
	ctx, cancel := context.WithTimeout(c, h.contextTimeout)
	defer cancel()

	ctxSession := ctx.Value(configs.ForContext)
	if ctxSession == nil {
		return models.User{}, models.HTTPError{
			Code:    http.StatusNotFound,
			Message: models.ErrContextNilError,
		}
	}
	currentSession, ok := ctxSession.(models.Session)
	if !ok {
		return models.User{}, models.HTTPError{
			Code:    http.StatusNotFound,
			Message: models.ErrConvertToSession,
		}
	}

	currentUser, err := h.UserRepo.GetUserByID(c, currentSession.UserID)
	if err != nil {
		return models.User{}, models.HTTPError{
			Code:    http.StatusNotFound,
			Message: err.Error(),
		}
	}

	err = currentUser.FillProfile(newUserData)
	if err != nil {
		return models.User{}, models.HTTPError{
			Code:    http.StatusBadRequest,
			Message: err.Error(),
		}
	}

	_, err = h.UserRepo.UpdateUser(c, currentUser)
	if err != nil {
		return models.User{}, models.HTTPError{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		}
	}

	return currentUser, models.StatusOk200
}

func (h *userUsecase) AddPhoto(c context.Context, photo io.Reader, fileName string) (models.Photo, models.HTTPError) {
	ctx, cancel := context.WithTimeout(c, h.contextTimeout)
	defer cancel()

	ctxSession := ctx.Value(configs.ForContext)
	if ctxSession == nil {
		return models.Photo{}, models.HTTPError{
			Code:    http.StatusNotFound,
			Message: models.ErrContextNilError,
		}
	}

	currentSession, ok := ctxSession.(models.Session)
	if !ok {
		return models.Photo{}, models.HTTPError{
			Code:    http.StatusNotFound,
			Message: models.ErrConvertToSession,
		}
	}

	user, err := h.UserRepo.GetUserByID(c, currentSession.UserID)
	if err != nil {
		return models.Photo{}, models.HTTPError{
			Code:    http.StatusNotFound,
			Message: err.Error(),
		}
	}

	photoPath, err := h.File.SaveUserPhoto(user, photo, fileName)
	if err != nil {
		return models.Photo{}, models.HTTPError{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		}
	}

	user.AddNewPhoto(photoPath)

	err = h.UserRepo.UpdateImgs(c, user.ID, user.Imgs)
	if err != nil {
		return models.Photo{}, models.HTTPError{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		}
	}

	return models.Photo{Path: photoPath}, models.StatusOk200
}

func (h *userUsecase) DeletePhoto(c context.Context, photo models.Photo) models.HTTPError {
	ctx, cancel := context.WithTimeout(c, h.contextTimeout)
	defer cancel()

	ctxSession := ctx.Value(configs.ForContext)
	if ctxSession == nil {
		return models.HTTPError{
			Code:    http.StatusNotFound,
			Message: models.ErrContextNilError,
		}
	}
	currentSession, ok := ctxSession.(models.Session)
	if !ok {
		return models.HTTPError{
			Code:    http.StatusNotFound,
			Message: models.ErrConvertToSession,
		}
	}

	user, err := h.UserRepo.GetUserByID(c, currentSession.UserID)
	if err != nil {
		return models.HTTPError{
			Code:    http.StatusNotFound,
			Message: err.Error(),
		}
	}

	err = user.DeletePhoto(photo)
	if err != nil {
		return models.HTTPError{
			Code:    http.StatusBadRequest,
			Message: err.Error(),
		}
	}

	err = h.UserRepo.UpdateImgs(c, user.ID, user.Imgs)
	if err != nil {
		return models.HTTPError{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		}
	}

	err = h.File.Delete(photo.Path)
	if err != nil {
		return models.HTTPError{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		}
	}

	return models.StatusOk200
}

// @Summary LogIn
// @Description log in
// @Tags login
// @Accept json
// @Produce json
// @Param input body LoginUser true "data for login"
// @Success 200 {object} JSON
// @Failure 400,404,500
// @Router /login [post]
func (h *userUsecase) Login(c context.Context, logUserData models.LoginUser) (models.User, models.HTTPError) {
	identifiableUser, err := h.UserRepo.GetUser(c, logUserData.Email)
	if err != nil {
		return models.User{}, models.HTTPError{
			Code:    http.StatusNotFound,
			Message: err.Error(),
		}
	}

	if !hasher.CheckWithHash(identifiableUser.Password, logUserData.Password) {
		return models.User{}, models.HTTPError{
			Code:    http.StatusNotFound,
			Message: "",
		}
	}

	return identifiableUser, models.StatusOk200
}

// @Summary SignUp
// @Description registration user
// @Tags registration
// @Accept json
// @Produce json
// @Param input body LoginUser true "data for registration"
// @Success 200 {object} JSON
// @Failure 400,404,500
// @Router /signup [post]
func (h *userUsecase) Signup(c context.Context, logUserData models.LoginUser) (models.User, models.HTTPError) {
	ctx, cancel := context.WithTimeout(c, h.contextTimeout)
	defer cancel()

	identifiableUser, _ := h.UserRepo.GetUser(ctx, logUserData.Email)
	if !identifiableUser.IsEmpty() {
		return models.User{}, models.HTTPError{
			Code:    models.StatusEmailAlreadyExists,
			Message: "",
		}
	}

	var err error
	logUserData.Password = hasher.HashAndSalt(nil, logUserData.Password)

	user, err := h.UserRepo.CreateUser(c, logUserData)
	if err != nil {
		return models.User{}, models.HTTPError{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		}
	}

	err = h.File.CreateFoldersForNewUser(user)
	if err != nil {
		return models.User{}, models.HTTPError{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		}
	}

	return user, models.StatusOk200
}

func (h *userUsecase) NextUser(c context.Context) ([]models.User, models.HTTPError) {
	ctx, cancel := context.WithTimeout(c, h.contextTimeout)
	defer cancel()

	ctxSession := ctx.Value(configs.ForContext)
	if ctxSession == nil {
		return nil, models.HTTPError{
			Code:    http.StatusNotFound,
			Message: models.ErrContextNilError,
		}
	}
	currentSession, ok := ctxSession.(models.Session)
	if !ok {
		return nil, models.HTTPError{
			Code:    http.StatusNotFound,
			Message: models.ErrConvertToSession,
		}
	}

	currentUser, err := h.UserRepo.GetUserByID(c, currentSession.UserID)
	if err != nil {
		return nil, models.HTTPError{
			Code:    http.StatusNotFound,
			Message: err.Error(),
		}
	}

	nextUsers, err := h.UserRepo.GetNextUserForSwipe(ctx, currentUser.ID)
	if err != nil {
		return nil, models.HTTPError{
			Code:    http.StatusNotFound,
			Message: models.ErrContextNilError,
		}
	}

	return nextUsers, models.StatusOk200
}

func (h *userUsecase) GetAllTags(c context.Context) (models.Tags, models.HTTPError) {
	ctx, cancel := context.WithTimeout(c, h.contextTimeout)
	defer cancel()

	ctxSession := ctx.Value(configs.ForContext)
	if ctxSession == nil {
		return models.Tags{}, models.HTTPError{
			Code:    http.StatusNotFound,
			Message: models.ErrContextNilError,
		}
	}
	currentSession, ok := ctxSession.(models.Session)
	if !ok {
		return models.Tags{}, models.HTTPError{
			Code:    http.StatusNotFound,
			Message: models.ErrConvertToSession,
		}
	}

	_, err := h.UserRepo.GetUserByID(c, currentSession.UserID)
	if err != nil {
		return models.Tags{}, models.HTTPError{
			Code:    http.StatusNotFound,
			Message: err.Error(),
		}
	}

	allTags, err := h.UserRepo.GetTags(ctx)
	if err != nil {
		return models.Tags{}, models.HTTPError{
			Code:    http.StatusNotFound,
			Message: err.Error(),
		}
	}
	var respTag models.Tag
	var currentAllTags = make(map[uint64]models.Tag)
	var respAllTags models.Tags
	counter := 0

	for _, value := range allTags {
		respTag.TagName = value
		currentAllTags[uint64(counter)] = respTag
		counter++
	}

	respAllTags.AllTags = currentAllTags
	respAllTags.Count = uint64(counter)

	return respAllTags, models.StatusOk200
}

func (h *userUsecase) UsersMatches(c context.Context) (models.Matches, models.HTTPError) {
	ctx, cancel := context.WithTimeout(c, h.contextTimeout)
	defer cancel()

	ctxSession := ctx.Value(configs.ForContext)
	if ctxSession == nil {
		return models.Matches{}, models.HTTPError{
			Code:    http.StatusNotFound,
			Message: models.ErrContextNilError,
		}
	}
	currentSession, ok := ctxSession.(models.Session)
	if !ok {
		return models.Matches{}, models.HTTPError{
			Code:    http.StatusNotFound,
			Message: models.ErrConvertToSession,
		}
	}

	_, err := h.UserRepo.GetUserByID(c, currentSession.UserID)
	if err != nil {
		return models.Matches{}, models.HTTPError{
			Code:    http.StatusNotFound,
			Message: err.Error(),
		}
	}

	// find matches
	mathes, err := h.UserRepo.GetUsersMatches(ctx, currentSession.UserID)
	if err != nil {
		return models.Matches{}, models.HTTPError{
			Code:    http.StatusNotFound,
			Message: err.Error(),
		}
	}

	// count
	counter := 0
	var allMathesMap = make(map[uint64]models.User)
	for _, value := range mathes {
		allMathesMap[uint64(counter)] = value
		counter++
	}

	var allMatches models.Matches
	allMatches.AllUsers = allMathesMap
	allMatches.Count = strconv.Itoa(counter)

	return allMatches, models.StatusOk200
}

func (h *userUsecase) Reaction(c context.Context, reactionData models.UserReaction) (models.Match, models.HTTPError) {
	ctx, cancel := context.WithTimeout(c, h.contextTimeout)
	defer cancel()

	ctxSession := ctx.Value(configs.ForContext)
	if ctxSession == nil {
		return models.Match{}, models.HTTPError{
			Code:    http.StatusNotFound,
			Message: models.ErrContextNilError,
		}
	}
	currentSession, ok := ctxSession.(models.Session)
	if !ok {
		return models.Match{}, models.HTTPError{
			Code:    http.StatusNotFound,
			Message: models.ErrConvertToSession,
		}
	}

	_, err := h.UserRepo.GetUserByID(c, currentSession.UserID)
	if err != nil {
		return models.Match{}, models.HTTPError{
			Code:    http.StatusNotFound,
			Message: err.Error(),
		}
	}

	// added reaction in db
	err = h.UserRepo.AddReaction(ctx, currentSession.UserID, reactionData.Id, reactionData.Reaction)
	if err != nil {
		return models.Match{}, models.HTTPError{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		}
	}

	// get users who liked current user
	var likes []uint64
	likes, err = h.UserRepo.GetLikes(ctx, currentSession.UserID)
	if err != nil {
		return models.Match{}, models.HTTPError{
			Code:    http.StatusNotFound,
			Message: err.Error(),
		}
	}

	var currMath models.Match
	currMath.Match = false
	for _, value := range likes {
		if value == reactionData.Id {
			currMath.Match = true
			err = h.UserRepo.DeleteLike(ctx, currentSession.UserID, reactionData.Id)
			if err != nil {
				return models.Match{}, models.HTTPError{
					Code:    http.StatusInternalServerError,
					Message: err.Error(),
				}
			}
			err = h.UserRepo.AddMatch(ctx, currentSession.UserID, reactionData.Id)
			if err != nil {
				return models.Match{}, models.HTTPError{
					Code:    http.StatusInternalServerError,
					Message: err.Error(),
				}
			}
		}
	}

	return currMath, models.StatusOk200
}
