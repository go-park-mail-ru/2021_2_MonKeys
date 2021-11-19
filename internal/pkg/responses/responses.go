package responses

import (
	"dripapp/internal/dripapp/models"
	"dripapp/internal/pkg/logger"
	"encoding/json"
	"io/ioutil"
	"net/http"
)

type JSON struct {
	Status int         `json:"status"`
	Body   interface{} `json:"body"`
}

func SendOK(w http.ResponseWriter) {
	resp := JSON{
		Status: http.StatusOK,
	}

	err := WriteJSON(w, resp)
	if err != nil {
		SendError(w, models.HTTPError{
			Code:    http.StatusInternalServerError,
			Message: models.ErrWriteByte,
		},
			logger.DripLogger.ErrorLogging,
		)
		return
	}

	logger.DripLogger.Info.Printf("CODE %d", resp.Status)
}

func SendData(w http.ResponseWriter, v interface{}) {
	resp := JSON{
		Status: http.StatusOK,
		Body: v,
	}

	err := WriteJSON(w, resp)
	if err != nil {
		SendError(w, models.HTTPError{
			Code:    http.StatusInternalServerError,
			Message: models.ErrWriteByte,
		},
			logger.DripLogger.ErrorLogging,
		)
		return
	}

	logger.DripLogger.Info.Printf("CODE %d", resp.Status)
}

func SendError(w http.ResponseWriter, httpErr models.HTTPError, logging func(int, string)) {
	resp := JSON{
		Status: httpErr.Code,
	}

	err := WriteJSON(w, resp)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	logging(httpErr.Code, httpErr.Message.Error())
}

func ReadJSON(r *http.Request, v interface{}) error {
	byteReq, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return err
	}

	err = json.Unmarshal(byteReq, v)
	if err != nil {
		return err
	}

	return nil
}

func WriteJSON(w http.ResponseWriter, v interface{}) error {
	byteResp, err := json.Marshal(v)
	if err != nil {
		return err
	}

	_, err = w.Write(byteResp)
	if err != nil {
		return err
	}

	return nil
}
