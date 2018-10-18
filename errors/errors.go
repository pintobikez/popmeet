package errors

import (
	"fmt"
	"gopkg.in/go-playground/validator.v9"
)

type ErrResponse struct {
	Error ErrContent `json:"error"`
}

type ErrContent struct {
	Code      int              `json:"code"`
	Msg       string           `json:"message"`
	ValErrors []*ErrValidation `json:"validation_errors,omitempty"`
}

type ErrValidation struct {
	Field string `json:"field"`
	Error string `json:"message"`
}

type HealthStatus struct {
	Repo *HealthStatusDetail `json:"repository"`
}

type HealthStatusDetail struct {
	Status string `json:"status"`
	Detail string `json:"detail,omitempty"`
}

const (
	ErrorInterestNotFound    = 1001
	ErrorInterestsNotFound   = 1002
	ErrorUserNotFound        = 1003
	ErrorUserProfileNotFound = 1004
	ErrorCreatingToken       = 1005
	ErrorEventNotFound       = 1006
	ErrorCantAddUSerToEvent  = 1007

	ValidationError = "Validation Errors"
	errorMessage    = "Field validation for %s failed on the '%s' tag"
)

// Processes the Validation errors into a map
func processValidationErrors(err error) []*ErrValidation {
	var ers []*ErrValidation
	for _, err := range err.(validator.ValidationErrors) {
		ers = append(ers, &ErrValidation{Field: err.Namespace(), Error: fmt.Sprintf(errorMessage, err.Field(), err.Tag())})
	}
	return ers
}

// Processes the Validation errors into a map
func GeneralErrorJson(code int, err string) *ErrResponse {
	return &ErrResponse{ErrContent{code, err, nil}}
}

func ValidationErrorJson(code int, err error) *ErrResponse {
	return &ErrResponse{ErrContent{code, ValidationError, processValidationErrors(err)}}
}
