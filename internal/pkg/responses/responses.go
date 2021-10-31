package responses

import (
	"dripapp/internal/pkg/models"
	"encoding/json"
	"log"
	"net/http"
)

func SendResp(resp models.JSON, w http.ResponseWriter) {
	byteResp, err := json.Marshal(resp)
	if err != nil {
		SendErrorResponse(w, &models.HTTPError{
			Code:    http.StatusInternalServerError,
			Message: "Error encoding json",
		})
		return
	}
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(byteResp)
	if err != nil {
		SendErrorResponse(w, &models.HTTPError{
			Code:    http.StatusInternalServerError,
			Message: "Error write byte",
		})
		return
	}
	log.Printf("CODE %d", resp.Status)
}

func SendErrorResponse(w http.ResponseWriter, error *models.HTTPError) {
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
	log.Printf("CODE %d ERROR %s", error.Code, error.Message)
}
