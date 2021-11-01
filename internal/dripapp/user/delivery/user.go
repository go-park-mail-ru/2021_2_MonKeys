package delivery

import (
	"dripapp/internal/dripapp/models"
	"dripapp/internal/pkg/responses"
	"encoding/json"
	"io/ioutil"
	"net/http"
)

const maxPhotoSize = 20 * 1024 * 1025

type UserHandler struct {
	// Logger    *zap.SugaredLogger
	SessionUcase models.SessionUsecase
	UserUCase    models.UserUsecase
}

func (h *UserHandler) CurrentUser(w http.ResponseWriter, r *http.Request) {
	var resp responses.JSON

	user, status := h.UserUCase.CurrentUser(r.Context())
	resp.Status = status.Code
	if status.Code != http.StatusOK {
		responses.SendErrorResponse(w, status)
		return
	}

	resp.Body = user
	responses.SendOKResp(resp, w)
}

func (h *UserHandler) EditProfileHandler(w http.ResponseWriter, r *http.Request) {
	var resp responses.JSON
	byteReq, err := ioutil.ReadAll(r.Body)
	if err != nil {
		responses.SendErrorResponse(w, models.HTTPError{
			Code:    http.StatusBadRequest,
			Message: err.Error(),
		})
		return
	}

	var newUserData models.User
	err = json.Unmarshal(byteReq, &newUserData)
	if err != nil {
		responses.SendErrorResponse(w, models.HTTPError{
			Code:    http.StatusBadRequest,
			Message: err.Error(),
		})
		return
	}

	user, status := h.UserUCase.EditProfile(r.Context(), newUserData)
	resp.Status = status.Code
	if status.Code != http.StatusOK {
		responses.SendErrorResponse(w, status)
	}

	resp.Body = user
	responses.SendOKResp(resp, w)
}

func (h *UserHandler) UploadPhoto(w http.ResponseWriter, r *http.Request) {
	var resp responses.JSON

	err := r.ParseMultipartForm(maxPhotoSize)
	if err != nil {
		responses.SendErrorResponse(w, models.HTTPError{
			Code:    http.StatusBadRequest,
			Message: err.Error(),
		})
		return
	}

	uploadedPhoto, fileHeader, err := r.FormFile("photo")
	if err != nil {
		responses.SendErrorResponse(w, models.HTTPError{
			Code:    http.StatusBadRequest,
			Message: err.Error(),
		})
		return
	}
	defer uploadedPhoto.Close()

	var photo models.Photo
	var status models.HTTPError

	photo.Path, status = h.UserUCase.AddPhoto(r.Context(), uploadedPhoto, fileHeader.Filename)
	resp.Status = status.Code
	if resp.Status != http.StatusOK {
		responses.SendErrorResponse(w, models.HTTPError{
			Code:    resp.Status,
			Message: err.Error(),
		})
		return
	}

	resp.Body = photo
	responses.SendOKResp(resp, w)
}

func (h *UserHandler) DeletePhoto(w http.ResponseWriter, r *http.Request) {
	var resp responses.JSON

	byteReq, err := ioutil.ReadAll(r.Body)
	if err != nil {
		responses.SendErrorResponse(w, models.HTTPError{
			Code:    http.StatusBadRequest,
			Message: err.Error(),
		})
		return
	}

	var photo models.Photo
	err = json.Unmarshal(byteReq, &photo)
	if err != nil {
		responses.SendErrorResponse(w, models.HTTPError{
			Code:    http.StatusBadRequest,
			Message: err.Error(),
		})
		return
	}

	status := h.UserUCase.DeletePhoto(r.Context(), photo)
	resp.Status = status.Code
	if status.Code != http.StatusOK {
		responses.SendErrorResponse(w, models.HTTPError{
			Code:    resp.Status,
			Message: err.Error(),
		})
		return
	}

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
// @Router /signup [post]
func (h *UserHandler) SignupHandler(w http.ResponseWriter, r *http.Request) {
	var resp responses.JSON

	byteReq, err := ioutil.ReadAll(r.Body)
	if err != nil {
		responses.SendErrorResponse(w, models.HTTPError{
			Code:    http.StatusBadRequest,
			Message: err.Error(),
		})
		return
	}

	var logUserData models.LoginUser
	err = json.Unmarshal(byteReq, &logUserData)
	if err != nil {
		resp.Status = http.StatusBadRequest
		responses.SendErrorResponse(w, models.HTTPError{
			Code:    http.StatusBadRequest,
			Message: err.Error(),
		})
		return
	}

	user, status := h.UserUCase.Signup(r.Context(), logUserData)
	resp.Status = status.Code
	if status.Code != http.StatusOK {
		responses.SendErrorResponse(w, status)
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
			Message: err.Error(),
		})
		return
	}
	resp.Body = user

	http.SetCookie(w, &cookie)

	responses.SendOKResp(resp, w)
}

func (h *UserHandler) NextUserHandler(w http.ResponseWriter, r *http.Request) {
	var resp responses.JSON

	nextUser, status := h.UserUCase.NextUser(r.Context())
	resp.Status = status.Code
	if status.Code != http.StatusOK {
		responses.SendErrorResponse(w, status)
	}

	resp.Body = nextUser
	responses.SendOKResp(resp, w)
}

func (h *UserHandler) GetAllTags(w http.ResponseWriter, r *http.Request) {
	var resp responses.JSON
	allTags, status := h.UserUCase.GetAllTags(r.Context())
	resp.Status = status.Code
	if status.Code != http.StatusOK {
		responses.SendErrorResponse(w, status)
	}

	resp.Body = allTags
	responses.SendOKResp(resp, w)
}

func (h *UserHandler) MatchesHandler(w http.ResponseWriter, r *http.Request) {
	var resp responses.JSON
	matches, status := h.UserUCase.UsersMatches(r.Context())
	resp.Status = status.Code
	if status.Code != http.StatusOK {
		responses.SendErrorResponse(w, status)
	}

	resp.Body = matches
	responses.SendOKResp(resp, w)
}

func (h *UserHandler) ReactionHandler(w http.ResponseWriter, r *http.Request) {
	var resp responses.JSON

	byteReq, err := ioutil.ReadAll(r.Body)
	if err != nil {
		responses.SendErrorResponse(w, models.HTTPError{
			Code:    http.StatusBadRequest,
			Message: err.Error(),
		})
		return
	}

	var reactionData models.UserReaction
	err = json.Unmarshal(byteReq, &reactionData)
	if err != nil {
		resp.Status = http.StatusBadRequest
		responses.SendErrorResponse(w, models.HTTPError{
			Code:    http.StatusBadRequest,
			Message: err.Error(),
		})
		return
	}

	match, status := h.UserUCase.Reaction(r.Context(), reactionData)

	resp.Status = status.Code
	if status.Code != http.StatusOK {
		responses.SendErrorResponse(w, status)
	}

	resp.Body = match
	responses.SendOKResp(resp, w)
}
