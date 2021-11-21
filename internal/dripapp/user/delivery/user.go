package delivery

import (
	"dripapp/internal/dripapp/models"
	"dripapp/internal/pkg/logger"
	"dripapp/internal/pkg/responses"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
)

const maxPhotoSize = 20 * 1024 * 1025

type UserHandler struct {
	SessionUcase models.SessionUsecase
	UserUCase    models.UserUsecase
	Logger       logger.Logger
}

func (h *UserHandler) CurrentUser(w http.ResponseWriter, r *http.Request) {
	user, err := h.UserUCase.CurrentUser(r.Context())
	if err != nil {
		responses.SendError(w, models.HTTPError{
			Code:    http.StatusNotFound,
			Message: err,
		}, h.Logger.ErrorLogging)
		return
	}

	responses.SendData(w, user)
}

func (h *UserHandler) EditProfileHandler(w http.ResponseWriter, r *http.Request) {
	var newUserData models.User
	err := responses.ReadJSON(r, &newUserData)
	if err != nil {
		responses.SendError(w, models.HTTPError{
			Code:    http.StatusBadRequest,
			Message: err,
		}, h.Logger.ErrorLogging)
		return
	}

	user, err := h.UserUCase.EditProfile(r.Context(), newUserData)
	if err != nil {
		responses.SendError(w, models.HTTPError{
			Code:    http.StatusNotFound,
			Message: err,
		}, h.Logger.ErrorLogging)
		return
	}

	responses.SendData(w, user)
}

func (h *UserHandler) UploadPhoto(w http.ResponseWriter, r *http.Request) {
	err := r.ParseMultipartForm(maxPhotoSize)
	if err != nil {
		responses.SendError(w, models.HTTPError{
			Code:    http.StatusBadRequest,
			Message: err,
		}, h.Logger.ErrorLogging)
		return
	}

	uploadedPhoto, fileHeader, err := r.FormFile("photo")
	if err != nil {
		responses.SendError(w, models.HTTPError{
			Code:    http.StatusBadRequest,
			Message: err,
		}, h.Logger.ErrorLogging)
		return
	}
	defer uploadedPhoto.Close()

	photo, err := h.UserUCase.AddPhoto(r.Context(), uploadedPhoto, fileHeader.Filename)
	if err != nil {
		responses.SendError(w, models.HTTPError{
			Code:    http.StatusInternalServerError,
			Message: err,
		}, h.Logger.ErrorLogging)
		return
	}

	responses.SendData(w, photo)
}

func (h *UserHandler) DeletePhoto(w http.ResponseWriter, r *http.Request) {
	var photo models.Photo
	err := responses.ReadJSON(r, &photo)
	if err != nil {
		responses.SendError(w, models.HTTPError{
			Code:    http.StatusBadRequest,
			Message: err,
		}, h.Logger.ErrorLogging)
		return
	}

	err = h.UserUCase.DeletePhoto(r.Context(), photo)
	if err != nil {
		responses.SendError(w, models.HTTPError{
			Code:    http.StatusNotFound,
			Message: err,
		}, h.Logger.ErrorLogging)
		return
	}

	responses.SendOK(w)
}

func (h *UserHandler) SignupHandler(w http.ResponseWriter, r *http.Request) {
	var logUserData models.LoginUser
	err := responses.ReadJSON(r, &logUserData)
	if err != nil {
		responses.SendError(w, models.HTTPError{
			Code:    http.StatusBadRequest,
			Message: err,
		}, h.Logger.ErrorLogging)
		return
	}

	user, err := h.UserUCase.Signup(r.Context(), logUserData)
	if err != nil {
		code := http.StatusNotFound
		if err == models.ErrEmailAlreadyExists {
			code = models.StatusEmailAlreadyExists
		}
		responses.SendError(w, models.HTTPError{
			Code:    code,
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
		responses.SendError(w, models.HTTPError{
			Code:    http.StatusInternalServerError,
			Message: err,
		}, h.Logger.WarnLogging)
		return
	}

	http.SetCookie(w, &cookie)

	responses.SendData(w, user)
}

func (h *UserHandler) NextUserHandler(w http.ResponseWriter, r *http.Request) {
	nextUser, err := h.UserUCase.NextUser(r.Context())
	if err != nil {
		responses.SendError(w, models.HTTPError{
			Code:    http.StatusNotFound,
			Message: err,
		}, h.Logger.ErrorLogging)
		return
	}

	responses.SendData(w, nextUser)
}

func (h *UserHandler) GetAllTags(w http.ResponseWriter, r *http.Request) {
	allTags, err := h.UserUCase.GetAllTags(r.Context())
	if err != nil {
		responses.SendError(w, models.HTTPError{
			Code:    http.StatusNotFound,
			Message: err,
		}, h.Logger.ErrorLogging)
		return
	}

	responses.SendData(w, allTags)
}

func (h *UserHandler) MatchesHandler(w http.ResponseWriter, r *http.Request) {
	matches, err := h.UserUCase.UsersMatches(r.Context())
	if err != nil {
		responses.SendError(w, models.HTTPError{
			Code:    http.StatusNotFound,
			Message: err,
		}, h.Logger.ErrorLogging)
		return
	}

	responses.SendData(w, matches)
}

func (h *UserHandler) SearchMatchesHandler(w http.ResponseWriter, r *http.Request) {
	var searchData models.Search
	err := responses.ReadJSON(r, &searchData)
	if err != nil {
		responses.SendError(w, models.HTTPError{
			Code:    http.StatusBadRequest,
			Message: err,
		}, h.Logger.ErrorLogging)
		return
	}

	matches, err := h.UserUCase.UsersMatchesWithSearching(r.Context(), searchData)
	if err != nil {
		responses.SendError(w, models.HTTPError{
			Code:    http.StatusNotFound,
			Message: err,
		}, h.Logger.ErrorLogging)
		return
	}

	responses.SendData(w, matches)
}

func (h *UserHandler) ReactionHandler(w http.ResponseWriter, r *http.Request) {
	var reactionData models.UserReaction
	err := responses.ReadJSON(r, &reactionData)
	if err != nil {
		responses.SendError(w, models.HTTPError{
			Code:    http.StatusBadRequest,
			Message: err,
		}, h.Logger.ErrorLogging)
		return
	}

	match, err := h.UserUCase.Reaction(r.Context(), reactionData)
	if err != nil {
		responses.SendError(w, models.HTTPError{
			Code:    http.StatusNotFound,
			Message: err,
		}, h.Logger.ErrorLogging)
		return
	}

	responses.SendData(w, match)
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func (h *UserHandler) Notifications(w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		status := models.HTTPError{
			Code:    http.StatusInternalServerError,
			Message: err,
		}
		responses.SendError(w, status, h.Logger.ErrorLogging)
		return
	}
	go h.sendNewMsgNotifications(ws)
}

func (h *UserHandler) sendNewMsgNotifications(client *websocket.Conn) {
	for {
		var msg models.Message

		err := client.ReadJSON(&msg)
		if err != nil {
			h.Logger.ErrorLogging(http.StatusBadRequest, "ReadJSON: "+err.Error())
			return
		}

		err = client.WriteJSON(msg)
		if err != nil {
			h.Logger.ErrorLogging(http.StatusBadRequest, "WriteJSON")
			return
		}
	}
}

func (h *UserHandler) LikesHandler(w http.ResponseWriter, r *http.Request) {
	likes, err := h.UserUCase.UserLikes(r.Context())
	if err != nil {
		responses.SendError(w, models.HTTPError{
			Code:    http.StatusNotFound,
			Message: err,
		}, h.Logger.ErrorLogging)
		return
	}

	responses.SendData(w, likes)
}

func (h *UserHandler) GetChat(w http.ResponseWriter, r *http.Request) {
	fromId, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		responses.SendError(w, models.HTTPError{
			Code:    http.StatusNotFound,
			Message: err,
		}, h.Logger.ErrorLogging)
		return
	}
	lastId, err := strconv.Atoi(mux.Vars(r)["lastId"])
	if err != nil {
		responses.SendError(w, models.HTTPError{
			Code:    http.StatusNotFound,
			Message: err,
		}, h.Logger.ErrorLogging)
		return
	}

	mses, err := h.UserUCase.GetChat(r.Context(), uint64(fromId), uint64(lastId))
	if err != nil {
		responses.SendError(w, models.HTTPError{
			Code:    http.StatusNotFound,
			Message: err,
		}, h.Logger.ErrorLogging)
		return
	}

	responses.SendData(w, mses)
}

func (h *UserHandler) GetChats(w http.ResponseWriter, r *http.Request) {
	chats, err := h.UserUCase.GetChats(r.Context())
	if err != nil {
		responses.SendError(w, models.HTTPError{
			Code:    http.StatusNotFound,
			Message: err,
		}, h.Logger.ErrorLogging)
		return
	}

	responses.SendData(w, chats)
}

func (h *UserHandler) SendMessage(w http.ResponseWriter, r *http.Request) {
	var ms models.Message
	err := responses.ReadJSON(r, &ms)
	if err != nil {
		responses.SendError(w, models.HTTPError{
			Code:    http.StatusBadRequest,
			Message: err,
		}, h.Logger.ErrorLogging)
		return
	}

	err = h.UserUCase.SendMessage(r.Context(), ms)
	if err != nil {
		responses.SendError(w, models.HTTPError{
			Code:    http.StatusNotFound,
			Message: err,
		}, h.Logger.ErrorLogging)
		return
	}

	responses.SendOK(w)
}
