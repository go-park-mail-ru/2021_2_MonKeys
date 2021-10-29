package delivery

import (
	"dripapp/internal/pkg/models"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
)

const (
	StatusOK                  = 200
	StatusBadRequest          = 400
	StatusNotFound            = 404
	StatusInternalServerError = 500
	StatusEmailAlreadyExists  = 1001
)

const maxPhotoSize = 20 * 1024 * 1025

type UserHandler struct {
	// Logger    *zap.SugaredLogger
	UserUCase models.UserUsecase
}

func sendResp(resp models.JSON, w http.ResponseWriter) {
	byteResp, err := json.Marshal(resp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(byteResp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (h *UserHandler) CurrentUser(w http.ResponseWriter, r *http.Request) {
	var resp models.JSON

	user, status := h.UserUCase.CurrentUser(r.Context(), r)
	resp.Status = status
	if status == StatusOK {
		resp.Body = user
	}

	sendResp(resp, w)
}

func (h *UserHandler) EditProfileHandler(w http.ResponseWriter, r *http.Request) {
	var resp models.JSON
	byteReq, err := ioutil.ReadAll(r.Body)
	if err != nil {
		resp.Status = StatusBadRequest
		sendResp(resp, w)
		log.Printf("CODE %d ERROR %s", resp.Status, err)
		return
	}

	var newUserData models.User
	err = json.Unmarshal(byteReq, &newUserData)
	if err != nil {
		resp.Status = StatusBadRequest
		sendResp(resp, w)
		log.Printf("CODE %d ERROR %s", resp.Status, err)
		return
	}

	user, status := h.UserUCase.EditProfile(r.Context(), newUserData, r)
	resp.Status = status
	if status == StatusOK {
		resp.Body = user
	}

	sendResp(resp, w)
}

func (h *UserHandler) UploadPhoto(w http.ResponseWriter, r *http.Request) {
	var resp models.JSON

	err := r.ParseMultipartForm(maxPhotoSize)
	if err != nil {
		resp.Status = StatusBadRequest
		sendResp(resp, w)
		log.Printf("CODE %d ERROR %s", resp.Status, err)
		return
	}

	uploadedPhoto, _, err := r.FormFile("photo")
	if err != nil {
		resp.Status = StatusBadRequest
		sendResp(resp, w)
		log.Printf("CODE %d ERROR %s", resp.Status, err)
		return
	}
	defer uploadedPhoto.Close()

	var photo models.Photo

	photo.Path, resp.Status = h.UserUCase.AddPhoto(r.Context(), uploadedPhoto, r)
	if resp.Status != StatusOK {
		log.Printf("CODE %d ERROR %s", resp.Status, err)
		sendResp(resp, w)
		return
	}

	resp.Body = photo
	sendResp(resp, w)
}

func (h *UserHandler) DeletePhoto(w http.ResponseWriter, r *http.Request) {
	var resp models.JSON

	byteReq, err := ioutil.ReadAll(r.Body)
	if err != nil {
		resp.Status = StatusBadRequest
		sendResp(resp, w)
		log.Printf("CODE %d ERROR %s", resp.Status, err)
		return
	}

	var photo models.Photo
	err = json.Unmarshal(byteReq, &photo)
	if err != nil {
		resp.Status = StatusBadRequest
		sendResp(resp, w)
		log.Printf("CODE %d ERROR %s", resp.Status, err)
		return
	}

	resp.Status = h.UserUCase.DeletePhoto(r.Context(), photo, r)
	if resp.Status != StatusOK {
		log.Printf("CODE %d ERROR %s", resp.Status, err)
		sendResp(resp, w)
		return
	}

	sendResp(resp, w)
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
	var resp models.JSON

	byteReq, err := ioutil.ReadAll(r.Body)
	if err != nil {
		resp.Status = StatusBadRequest
		sendResp(resp, w)
		log.Printf("CODE %d ERROR %s", resp.Status, err)
		return
	}

	var logUserData models.LoginUser
	err = json.Unmarshal(byteReq, &logUserData)
	if err != nil {
		resp.Status = StatusBadRequest
		sendResp(resp, w)
		log.Printf("CODE %d ERROR %s", resp.Status, err)
		return
	}

	status := h.UserUCase.Signup(r.Context(), logUserData, w)
	resp.Status = status
	sendResp(resp, w)
}

func (h *UserHandler) NextUserHandler(w http.ResponseWriter, r *http.Request) {
	var resp models.JSON

	// get swiped usedata for registrationr id from json
	// byteReq, err := ioutil.ReadAll(r.Body)
	// if err != nil {
	// 	resp.Status = StatusBadRequest
	// 	sendResp(resp, w)
	// 	log.Printf("CODE %d ERROR %s", resp.Status, err)
	// 	return
	// }
	// var swipedUserData models.SwipedUser
	// var byteReq []byte
	// err := json.Unmarshal(byteReq, &swipedUserData)
	// if err != nil {
	// 	resp.Status = StatusBadRequest
	// 	sendResp(resp, w)
	// 	log.Printf("CODE %d ERROR %s", resp.Status, err)
	// 	return
	// }
	nextUser, status := h.UserUCase.NextUser(r.Context(), r)
	resp.Status = status
	if status == StatusOK {
		resp.Body = nextUser
	}

	sendResp(resp, w)
}

func (h *UserHandler) GetAllTags(w http.ResponseWriter, r *http.Request) {
	var resp models.JSON
	allTags, status := h.UserUCase.GetAllTags(r.Context(), r)
	resp.Body = allTags
	resp.Status = status
	sendResp(resp, w)
}
