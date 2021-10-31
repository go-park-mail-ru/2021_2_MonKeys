package models

type HTTPError struct {
	Code    int    `json:"-"`
	Message string `json:"error_description"`
}

var (
	InternalServerError500 = HTTPError{500, "InternalServerError500"}
	StatusOk200            = HTTPError{200, ""}

	StatusEmailAlreadyExists = 1001

	ErrContextNilError  = "context nil error"
	ErrConvertToSession = "convert to model session error"
	ErrExtractContext   = "context extract error"
	ErrAuth             = "auth error"
)
