package delivery

import (
	"dripapp/internal/dripapp/models"
	"dripapp/internal/pkg/logger"
	"dripapp/internal/pkg/responses"
	"encoding/json"
	"io/ioutil"
	"net/http"
)

const maxPhotoSize = 20 * 1024 * 1025

type UserHandler struct {
	SessionUcase models.SessionUsecase
	UserUCase    models.UserUsecase
	Logger       logger.Logger
}

func (h *UserHandler) CurrentUser(w http.ResponseWriter, r *http.Request) {
	var resp responses.JSON

	user, status := h.UserUCase.CurrentUser(r.Context())
	if status != nil {
		resp.Status = http.StatusNotFound
		responses.SendErrorResponse(w, models.HTTPError{
			Code:    http.StatusNotFound,
			Message: status,
		}, h.Logger.ErrorLogging)
		return
	}

	resp.Status = http.StatusOK
	resp.Body = user
	responses.SendOKResp(resp, w)
}

func (h *UserHandler) EditProfileHandler(w http.ResponseWriter, r *http.Request) {
	var resp responses.JSON
	byteReq, err := ioutil.ReadAll(r.Body)
	if err != nil {
		resp.Status = http.StatusBadRequest
		responses.SendErrorResponse(w, models.HTTPError{
			Code:    http.StatusBadRequest,
			Message: err,
		}, h.Logger.ErrorLogging)
		return
	}

	var newUserData models.User
	err = json.Unmarshal(byteReq, &newUserData)
	if err != nil {
		resp.Status = http.StatusBadRequest
		responses.SendErrorResponse(w, models.HTTPError{
			Code:    http.StatusBadRequest,
			Message: err,
		}, h.Logger.ErrorLogging)
		return
	}

	user, err := h.UserUCase.EditProfile(r.Context(), newUserData)
	if err != nil {
		responses.SendErrorResponse(w, models.HTTPError{
			Code:    http.StatusNotFound,
			Message: err,
		}, h.Logger.ErrorLogging)
		return
	}

	resp.Status = http.StatusOK
	resp.Body = user
	responses.SendOKResp(resp, w)
}

func (h *UserHandler) UploadPhoto(w http.ResponseWriter, r *http.Request) {
	var resp responses.JSON

	err := r.ParseMultipartForm(maxPhotoSize)
	if err != nil {
		responses.SendErrorResponse(w, models.HTTPError{
			Code:    http.StatusBadRequest,
			Message: err,
		}, h.Logger.ErrorLogging)
		return
	}

	uploadedPhoto, fileHeader, err := r.FormFile("photo")
	if err != nil {
		responses.SendErrorResponse(w, models.HTTPError{
			Code:    http.StatusBadRequest,
			Message: err,
		}, h.Logger.ErrorLogging)
		return
	}
	defer uploadedPhoto.Close()

	photo, status := h.UserUCase.AddPhoto(r.Context(), uploadedPhoto, fileHeader.Filename)
	if status != nil {
		resp.Status = http.StatusNotFound
		responses.SendErrorResponse(w, models.HTTPError{
			Code:    resp.Status,
			Message: status,
		}, h.Logger.ErrorLogging)
		return
	}

	resp.Status = http.StatusOK
	resp.Body = photo
	responses.SendOKResp(resp, w)
}

func (h *UserHandler) DeletePhoto(w http.ResponseWriter, r *http.Request) {
	var resp responses.JSON

	byteReq, err := ioutil.ReadAll(r.Body)
	if err != nil {
		responses.SendErrorResponse(w, models.HTTPError{
			Code:    http.StatusBadRequest,
			Message: err,
		}, h.Logger.ErrorLogging)
		return
	}

	var photo models.Photo
	err = json.Unmarshal(byteReq, &photo)
	if err != nil {
		responses.SendErrorResponse(w, models.HTTPError{
			Code:    http.StatusBadRequest,
			Message: err,
		}, h.Logger.ErrorLogging)
		return
	}

	status := h.UserUCase.DeletePhoto(r.Context(), photo)
	if status != nil {
		resp.Status = http.StatusNotFound
		responses.SendErrorResponse(w, models.HTTPError{
			Code:    resp.Status,
			Message: status,
		}, h.Logger.ErrorLogging)
		return
	}

	resp.Status = http.StatusOK
	responses.SendOKResp(resp, w)
}

// @Summary SignUp
// @Description registration user
// @Tags registration
// @Accept json
// @Produce json
// @Param input body LoginUser true "data for registration"
// @Success 200 {object} JSON
// @Failure 400,404,500
// @Router /profile [post]
func (h *UserHandler) SignupHandler(w http.ResponseWriter, r *http.Request) {
	var resp responses.JSON

	byteReq, err := ioutil.ReadAll(r.Body)
	if err != nil {
		responses.SendErrorResponse(w, models.HTTPError{
			Code:    http.StatusBadRequest,
			Message: err,
		}, h.Logger.ErrorLogging)
		return
	}

	var logUserData models.LoginUser
	err = json.Unmarshal(byteReq, &logUserData)
	if err != nil {
		resp.Status = http.StatusBadRequest
		responses.SendErrorResponse(w, models.HTTPError{
			Code:    http.StatusBadRequest,
			Message: err,
		}, h.Logger.ErrorLogging)
		return
	}

	user, err := h.UserUCase.Signup(r.Context(), logUserData)
	if err != nil {
		resp.Status = http.StatusNotFound
		if err == models.ErrEmailAlreadyExists {
			resp.Status = models.StatusEmailAlreadyExists
		}
		responses.SendErrorResponse(w, models.HTTPError{
			Code:    resp.Status,
			Message: err,
		}, h.Logger.ErrorLogging)
		return
	}
	cookie := models.CreateSessionCookie(logUserData)

	sess := models.Session{
		Cookie: cookie.Value,
		UserID: user.ID,
	}
	err = h.SessionUcase.AddSession(r.Context(), sess)
	if err != nil {
		responses.SendErrorResponse(w, models.HTTPError{
			Code:    http.StatusInternalServerError,
			Message: err,
		}, h.Logger.WarnLogging)
		return
	}

	resp.Status = http.StatusOK
	resp.Body = user

	http.SetCookie(w, &cookie)

	responses.SendOKResp(resp, w)
}

func (h *UserHandler) NextUserHandler(w http.ResponseWriter, r *http.Request) {
	var resp responses.JSON

	nextUser, status := h.UserUCase.NextUser(r.Context())
	if status != nil {
		resp.Status = http.StatusNotFound
		responses.SendErrorResponse(w, models.HTTPError{
			Code:    resp.Status,
			Message: status,
		}, h.Logger.ErrorLogging)
		return
	}

	resp.Status = http.StatusOK
	resp.Body = nextUser
	responses.SendOKResp(resp, w)
}

func (h *UserHandler) GetAllTags(w http.ResponseWriter, r *http.Request) {
	var resp responses.JSON
	allTags, status := h.UserUCase.GetAllTags(r.Context())
	if status != nil {
		resp.Status = http.StatusNotFound
		responses.SendErrorResponse(w, models.HTTPError{
			Code:    resp.Status,
			Message: status,
		}, h.Logger.ErrorLogging)
		return
	}

	resp.Status = http.StatusOK
	resp.Body = allTags
	responses.SendOKResp(resp, w)
}

func (h *UserHandler) MatchesHandler(w http.ResponseWriter, r *http.Request) {
	var resp responses.JSON
	matches, status := h.UserUCase.UsersMatches(r.Context())
	if status != nil {
		resp.Status = http.StatusNotFound
		responses.SendErrorResponse(w, models.HTTPError{
			Code:    resp.Status,
			Message: status,
		}, h.Logger.ErrorLogging)
		return
	}

	resp.Status = http.StatusOK
	resp.Body = matches
	responses.SendOKResp(resp, w)
}

func (h *UserHandler) ReactionHandler(w http.ResponseWriter, r *http.Request) {
	var resp responses.JSON

	byteReq, err := ioutil.ReadAll(r.Body)
	if err != nil {
		responses.SendErrorResponse(w, models.HTTPError{
			Code:    http.StatusBadRequest,
			Message: err,
		}, h.Logger.ErrorLogging)
		return
	}

	var reactionData models.UserReaction
	err = json.Unmarshal(byteReq, &reactionData)
	if err != nil {
		resp.Status = http.StatusBadRequest
		responses.SendErrorResponse(w, models.HTTPError{
			Code:    http.StatusBadRequest,
			Message: err,
		}, h.Logger.ErrorLogging)
		return
	}

	match, status := h.UserUCase.Reaction(r.Context(), reactionData)

	if status != nil {
		resp.Status = http.StatusNotFound
		responses.SendErrorResponse(w, models.HTTPError{
			Code:    resp.Status,
			Message: status,
		}, h.Logger.ErrorLogging)
		return
	}

	resp.Status = http.StatusOK
	resp.Body = match
	responses.SendOKResp(resp, w)
}
