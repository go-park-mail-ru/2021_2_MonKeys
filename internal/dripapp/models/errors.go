package models

import "errors"

type HTTPError struct {
	Code    int    `json:"-"`
	Message error `json:"error_description"`
}

var (
	InternalServerError500 = HTTPError{500, errors.New("InternalServerError500")}
	StatusOk200            = HTTPError{200, errors.New("")}

	ErrNoSuchPhoto = errors.New("user does not have such a photo")

	StatusEmailAlreadyExists = 1001

	ErrContextNilError    = errors.New("context nil error")
	ErrConvertToSession   = errors.New("convert to model session error")
	ErrConvertToUser      = errors.New("convert to model user error")
	ErrExtractContext     = errors.New("context extract error")
	ErrAuth               = errors.New("auth error")
	ErrEmailAlreadyExists = errors.New("email already exists")

	ErrNoPermission = errors.New("no permission")
	ErrCSRF = errors.New("csrf-protection")
	ErrJson = errors.New("Error encoding json")
	ErrWriteByte = errors.New("Error write byte")

	ErrSessionAlreadyExists = errors.New("session already exists")
)
