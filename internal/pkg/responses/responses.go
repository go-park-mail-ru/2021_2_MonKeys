package responses

import (
	"dripapp/internal/dripapp/models"
	"encoding/json"
	"log"
	"net/http"
)

type JSON struct {
	Status int         `json:"status"`
	Body   interface{} `json:"body"`
}

func SendOKResp(resp JSON, w http.ResponseWriter) {
	byteResp, err := json.Marshal(resp)
	if err != nil {
		SendErrorResponse(w, models.HTTPError{
			Code:    http.StatusInternalServerError,
			Message: "Error encoding json",
		})
		return
	}
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(byteResp)
	if err != nil {
		SendErrorResponse(w, models.HTTPError{
			Code:    http.StatusInternalServerError,
			Message: "Error write byte",
		})
		return
	}
	log.Printf("CODE %d", resp.Status)
}

func SendErrorResponse(w http.ResponseWriter, httpErr models.HTTPError) {
	w.WriteHeader(httpErr.Code)
	body, err := json.Marshal(httpErr)
	if err != nil {
		SendErrorResponse(w, models.HTTPError{
			Code:    http.StatusInternalServerError,
			Message: "Error encoding json",
		})
		return
	}
	_, _ = w.Write(body)
	log.Printf("CODE %d ERROR %s", httpErr.Code, httpErr.Message)
}
