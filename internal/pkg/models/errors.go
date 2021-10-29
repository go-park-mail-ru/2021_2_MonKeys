package models

type HTTPError struct {
	Code    int    `json:"-"`
	Message string `json:"error_description"`
}

var (
	InternalServerError500 = HTTPError{500, "InternalServerError500"}
)
