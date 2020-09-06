package errs

import (
	"errors"
	"fmt"
	"net/http"
)

// Error constants for 400 bad request
const (
	ErrParameterRequired  = "E400001"
	ErrMaxImages          = "E400002"
	ErrMaxLimit           = "E400003"
	ErrRequestBodyInvalid = "E400004"
	ErrAlbumNotExist      = "E400005"
	ErrImageExist         = "E400006"
	ErrAlbumExist         = "E400007"
	ErrImageIDNotFound    = "E400008"
)

// Error constants for 404 bad request
const (
	ErrResourceNotFound = "E404001"
	ErrCodeNotFound     = "E404002"
)

// Error constants for 500 bad request
const (
	ErrInternalAppError = "E500002"
	ErrInternalDBError  = "E500003"
)

// Error constants for 504 Gateway timeout
const (
	ErrGatewayTimeout = "E504001"
)

// Error constants for 422 Unprocessable Entity
const (
	ErrEmptyBodyContent = "E422001"
)

// Errors - maps of error with the error code
type Errors map[int]map[string]string

// Error - return error response
type Error struct {
	Error    string `json:"error"`
	HTTPCode int    `json:"http_code,string"`
}

// ErrorResponse - Use to trow the errors to users
type ErrorResponse struct {
	Error string `json:"error"`
}

var errs *Errors
var allErrs map[string]Error

// Init function for errs package
func init() {
	errs = &Errors{
		http.StatusBadRequest: {
			ErrParameterRequired: "Parameter `%s` is a required field",
			ErrMaxImages:         "max no of images for an album reached",
			ErrMaxLimit:          "max only one image uplaod allowed",
			ErrAlbumNotExist:     "Album doesn't exist",
			ErrImageExist:        "Image already exist",
			ErrAlbumExist:        "Album already exist",
			ErrImageIDNotFound:   "Image id not found",
		},
		http.StatusNotFound: {
			ErrResourceNotFound: "Resource Not found",
			ErrCodeNotFound:     "Error code not found",
		},
		http.StatusInternalServerError: {
			ErrInternalAppError: "Internal Application Error, `%s`",
			ErrInternalDBError:  "Database Error, `%s`",
		},
		http.StatusGatewayTimeout: {
			ErrGatewayTimeout: "Gateway Timeout",
		},
		http.StatusUnprocessableEntity: {
			ErrEmptyBodyContent: "Cannot parse empty body",
		},
	}
	allErrs = make(map[string]Error)

	for httpcode, err := range *errs {
		for code, msg := range err {
			tmp := &Error{HTTPCode: httpcode, Error: msg}
			allErrs[code] = *tmp
		}
	}
}

// GetErrorByCode ...
func GetErrorByCode(code string) (res Error, err error) {
	var ok bool
	if res, ok = allErrs[code]; !ok {
		err = errors.New(ErrCodeNotFound)
		return
	}
	return
}

// GetErrors ...
func GetErrors() (res map[string]Error) {
	res = allErrs
	return
}

// FormateErrorResponse ...
func FormateErrorResponse(mErr Error, val ...interface{}) (res ErrorResponse) {
	if len(val) > 0 {
		mErr.Error = fmt.Sprintf(mErr.Error, val...)
	}

	errRes := &ErrorResponse{
		Error: mErr.Error,
	}
	return *errRes
}
