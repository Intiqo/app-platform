package domain

import "fmt"

type NotFoundError struct{}

func (e NotFoundError) Error() string {
	return "The resource you are looking for does not exist"
}

// InvalidRequestError defines model for invalid request error.
type InvalidRequestError struct {
	Message string `json:"message" example:"invalid request"`
}

func (e InvalidRequestError) Error() string {
	return e.Message
}

// UnauthorizedError defines model for unauthorized error.
type UnauthorizedError struct {
	Code    string `json:"code" example:"UNAUTHORIZED"`
	Message string `json:"message" example:"You are not authorized to access this resource"`
}

func (e UnauthorizedError) Error() string {
	return e.Message
}

// ForbiddenAccessError defines model for forbidden access error.
type ForbiddenAccessError struct {
	Code    string `json:"code" example:"FORBIDDEN_ACCESS"`
	Message string `json:"message" example:"You are forbidden from accessing this resource"`
}

func (e ForbiddenAccessError) Error() string {
	return e.Message
}

type ValidationError struct {
	Code    string   `json:"code" example:"VALIDATION_ERROR"`
	Message string   `json:"message" example:"Not a valid mobile number"`
	Fields  []string `json:"fields" example:"mobile_number is required"`
}

func (e ValidationError) Error() string {
	if len(e.Fields) > 0 {
		return fmt.Sprintf(e.Message, e.Fields)
	}
	return e.Message
}

type UserError struct {
	Code    string `json:"code" example:"INVALID_REQUEST"`
	Message string `json:"message" example:"Oops! Something went wrong. Please try again later"`
}

func (e UserError) Error() string {
	return e.Message
}

type DataNotFoundError struct{}

func (e DataNotFoundError) Error() string {
	return "The record you are looking for does not exist"
}

type SystemError struct {
	Code    string `json:"code" example:"INTERNAL_SERVER_ERROR"`
	Message string `json:"message" example:"Oops! Something went wrong. Please try again later"`
}

func (e SystemError) Error() string {
	return e.Message
}

const (
	ErrorCodeINVALIDREQUEST      = "INVALID_REQUEST"
	ErrorCodeVALIDATIONERROR     = "VALIDATION_ERROR"
	ErrorCodeINTERNALSERVERERROR = "INTERNAL_SERVER_ERROR"
	ErrorCodeUNAUTHORIZED        = "UNAUTHORIZED"
	ErrorCodeFORBIDDENACCESS     = "FORBIDDEN_ACCESS"
)
