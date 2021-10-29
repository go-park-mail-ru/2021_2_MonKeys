package responses

import (
	"dripapp/internal/pkg/models"
	"encoding/json"
	"net/http"

	"github.com/sirupsen/logrus"
)

func SendResp(resp models.JSON, w http.ResponseWriter) {
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

func SendErrorResponse(w http.ResponseWriter, error *models.HTTPError) {
	logrus.Error(error.Message)
	w.WriteHeader(error.Code)
	body, err := json.Marshal(error)
	if err != nil {
		SendErrorResponse(w, &models.HTTPError{
			Code:    http.StatusInternalServerError,
			Message: "Error encoding json",
		})
		return
	}
	_, _ = w.Write(body)
}
