package delivery

import (
	"dripapp/internal/pkg/models"
	"dripapp/internal/pkg/responses"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

type UserHandler struct {
	// Logger    *zap.SugaredLogger
	SessionUcase models.SessionUsecase
	UserUCase    models.UserUsecase
}

func (h *UserHandler) CurrentUser(w http.ResponseWriter, r *http.Request) {
	var resp models.JSON

	user, status := h.UserUCase.CurrentUser(r.Context())
	resp.Status = status
	if status == http.StatusOK {
		resp.Body = user
	}

	responses.SendResp(resp, w)
}

func (h *UserHandler) EditProfileHandler(w http.ResponseWriter, r *http.Request) {
	var resp models.JSON
	byteReq, err := ioutil.ReadAll(r.Body)
	if err != nil {
		resp.Status = http.StatusBadRequest
		responses.SendResp(resp, w)
		log.Printf("CODE %d ERROR %s", resp.Status, err)
		return
	}

	var newUserData models.User
	err = json.Unmarshal(byteReq, &newUserData)
	if err != nil {
		resp.Status = http.StatusBadRequest
		responses.SendResp(resp, w)
		log.Printf("CODE %d ERROR %s", resp.Status, err)
		return
	}

	user, status := h.UserUCase.EditProfile(r.Context(), newUserData)
	resp.Status = status
	if status == http.StatusOK {
		resp.Body = user
	}

	responses.SendResp(resp, w)
}

func (h *UserHandler) UploadPhoto(w http.ResponseWriter, r *http.Request) {
	h.UserUCase.AddPhoto(r.Context(), w, r)
}

func (h *UserHandler) DeletePhoto(w http.ResponseWriter, r *http.Request) {
	h.UserUCase.DeletePhoto(r.Context(), w, r)
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
		resp.Status = http.StatusBadRequest
		responses.SendResp(resp, w)
		log.Printf("CODE %d ERROR %s", resp.Status, err)
		return
	}

	var logUserData models.LoginUser
	err = json.Unmarshal(byteReq, &logUserData)
	if err != nil {
		resp.Status = http.StatusBadRequest
		// responses.SendResp(resp, w)
		// log.Printf("CODE %d ERROR %s", resp.Status, err)
		responses.SendErrorResponse(w, &models.HTTPError{
			Code:    http.StatusBadRequest,
			Message: fmt.Sprintf("%v", err),
		})
		return
	}

	log.Println("Email: ", logUserData.Email, " Password: ", logUserData.Password)
	user, status := h.UserUCase.Signup(r.Context(), logUserData)
	resp.Status = status
	if status == http.StatusOK {
		cookie := models.CreateSessionCookie(logUserData)

		sess := models.Session{
			Cookie: cookie.Value,
			UserID: user.ID,
		}
		err = h.SessionUcase.AddSession(r.Context(), sess)
		if err != nil {
			// resp.Status = StatusInternalServerError
			// log.Printf("CODE %d ERROR %s", resp.Status, err)
			// responses.SendResp(resp, w)
			responses.SendErrorResponse(w, &models.HTTPError{
				Code:    http.StatusInternalServerError,
				Message: fmt.Sprintf("%v", err),
			})
			return
		}
		resp.Body = user

		http.SetCookie(w, &cookie)
	}
	responses.SendResp(resp, w)
}

func (h *UserHandler) NextUserHandler(w http.ResponseWriter, r *http.Request) {
	var resp models.JSON

	// get swiped usedata for registrationr id from json
	// byteReq, err := ioutil.ReadAll(r.Body)
	// if err != nil {
	// 	resp.Status = http.StatusBadRequest
	// 	responses.SendResp(resp, w)
	// 	log.Printf("CODE %d ERROR %s", resp.Status, err)
	// 	return
	// }
	// var swipedUserData models.SwipedUser
	// var byteReq []byte
	// err := json.Unmarshal(byteReq, &swipedUserData)
	// if err != nil {
	// 	resp.Status = http.StatusBadRequest
	// 	responses.SendResp(resp, w)
	// 	log.Printf("CODE %d ERROR %s", resp.Status, err)
	// 	return
	// }
	nextUser, status := h.UserUCase.NextUser(r.Context())
	resp.Status = status
	if status == http.StatusOK {
		resp.Body = nextUser
	}

	responses.SendResp(resp, w)
}

func (h *UserHandler) GetAllTags(w http.ResponseWriter, r *http.Request) {
	var resp models.JSON
	allTags, status := h.UserUCase.GetAllTags(r.Context())
	resp.Body = allTags
	resp.Status = status
	responses.SendResp(resp, w)
}
