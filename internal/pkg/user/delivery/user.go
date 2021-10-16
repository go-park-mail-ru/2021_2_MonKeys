package delivery

import (
	"dripapp/internal/pkg/models"
	"net/http"
)

const (
	StatusOK                  = 200
	StatusBadRequest          = 400
	StatusNotFound            = 404
	StatusInternalServerError = 500
	StatusEmailAlreadyExists  = 1001
)

type UserHandler struct {
	// Logger    *zap.SugaredLogger
	UserUCase models.UserUsecase
}

func (h *UserHandler) CurrentUser(w http.ResponseWriter, r *http.Request) {
	h.UserUCase.CurrentUser(w, r)
}

func (h *UserHandler) EditProfileHandler(w http.ResponseWriter, r *http.Request) {
	h.UserUCase.EditProfileHandler(w, r)
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
func (h *UserHandler) LoginHandler(w http.ResponseWriter, r *http.Request) {
	h.UserUCase.LoginHandler(w, r)
}

func (h *UserHandler) LogoutHandler(w http.ResponseWriter, r *http.Request) {
	h.UserUCase.LogoutHandler(w, r)
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
func (h *UserHandler) SignupHandler(w http.ResponseWriter, r *http.Request) {
	h.UserUCase.SignupHandler(w, r)
}

func (h *UserHandler) NextUserHandler(w http.ResponseWriter, r *http.Request) {
	h.UserUCase.NextUserHandler(w, r)
}

func (h *UserHandler) GetAllTags(w http.ResponseWriter, r *http.Request) {
	h.UserUCase.GetAllTags(w, r)
}