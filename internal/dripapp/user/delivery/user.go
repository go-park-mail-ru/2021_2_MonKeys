package delivery

import (
	"dripapp/internal/dripapp/models"
	_sessionModels "dripapp/internal/microservices/auth/models"
	"dripapp/internal/pkg/logger"
	"dripapp/internal/pkg/responses"
	"net/http"
)

const maxPhotoSize = 20 * 1024 * 1025

type UserHandler struct {
	SessionUcase _sessionModels.SessionUsecase
	UserUCase    models.UserUsecase
	Logger       logger.Logger
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

func (h *UserHandler) GetAllReports(w http.ResponseWriter, r *http.Request) {
	allReports, err := h.UserUCase.GetAllReports(r.Context())
	if err != nil {
		responses.SendError(w, models.HTTPError{
			Code:    http.StatusNotFound,
			Message: err,
		}, h.Logger.ErrorLogging)
		return
	}

	responses.SendData(w, allReports)
}

func (h *UserHandler) AddReport(w http.ResponseWriter, r *http.Request) {
	var reportData models.NewReport
	err := responses.ReadJSON(r, &reportData)
	if err != nil {
		responses.SendError(w, models.HTTPError{
			Code:    http.StatusBadRequest,
			Message: err,
		}, h.Logger.ErrorLogging)
		return
	}

	err = h.UserUCase.AddReport(r.Context(), reportData)
	if err != nil {
		responses.SendError(w, models.HTTPError{
			Code:    http.StatusNotFound,
			Message: err,
		}, h.Logger.ErrorLogging)
		return
	}

	responses.SendOK(w)
}

// func (h *UserHandler) UpdatePayment(w http.ResponseWriter, r *http.Request) {
// 	paymentId, err := strconv.Atoi(mux.Vars(r)["id"])
// 	if err != nil {
// 		responses.SendError(w, models.HTTPError{
// 			Code:    http.StatusNotFound,
// 			Message: err,
// 		}, h.Logger.ErrorLogging)
// 		return
// 	}

// 	err = h.UserUCase.UpdatePayment(r.Context(), uint64(paymentId))
// 	if err != nil {
// 		responses.SendError(w, models.HTTPError{
// 			Code:    http.StatusNotFound,
// 			Message: err,
// 		}, h.Logger.ErrorLogging)
// 		return
// 	}

// 	responses.SendOK(w)
// }

func (h *UserHandler) CreatePayment(w http.ResponseWriter, r *http.Request) {
	var paymentData models.Payment
	err := responses.ReadJSON(r, &paymentData)
	if err != nil {
		responses.SendError(w, models.HTTPError{
			Code:    http.StatusBadRequest,
			Message: err,
		}, h.Logger.ErrorLogging)
		return
	}

	redirectUrl, err := h.UserUCase.CreatePayment(r.Context(), paymentData)
	if err != nil {
		responses.SendError(w, models.HTTPError{
			Code:    http.StatusNotFound,
			Message: err,
		}, h.Logger.ErrorLogging)
		return
	}

	responses.SendData(w, redirectUrl)
}

func (h *UserHandler) HandlePaymentNotification(w http.ResponseWriter, r *http.Request) {
	var paymentNotificationData models.PaymentNotification
	err := responses.ReadJSON(r, &paymentNotificationData)
	if err != nil {
		responses.SendError(w, models.HTTPError{
			Code:    http.StatusBadRequest,
			Message: err,
		}, h.Logger.ErrorLogging)
		return
	}

	err = h.UserUCase.UpdatePayment(r.Context(), paymentNotificationData)
	if err != nil {
		responses.SendError(w, models.HTTPError{
			Code:    http.StatusNotFound,
			Message: err,
		}, h.Logger.ErrorLogging)
		return
	}

	responses.SendOK(w)
}

func (h *UserHandler) CheckSubscription(w http.ResponseWriter, r *http.Request) {
	payment, err := h.UserUCase.CheckSubscription(r.Context())
	if err != nil {
		responses.SendError(w, models.HTTPError{
			Code:    http.StatusNotFound,
			Message: err,
		}, h.Logger.ErrorLogging)
		return
	}

	responses.SendData(w, payment)
}
