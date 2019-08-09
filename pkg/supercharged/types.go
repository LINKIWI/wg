package supercharged

import (
	"encoding/json"
	"fmt"
)

const (
	// CodeServerUndefined describes an undefined server-side error.
	CodeServerUndefined = -1
	// CodeClientUndefined describes an undefined client-side error.
	CodeClientUndefined = -2
	// CodeInvalidParameters indicates the client supplied invalid input parameters.
	CodeInvalidParameters = -3
	// CodeNotFound indicates the client attempted to access an unknown route.
	CodeNotFound = -4
)

// Response formalizes a canonical Supercharged JSON response body.
type Response struct {
	Success bool            `json:"success"`
	Code    int             `json:"code"`
	Message string          `json:"message"`
	Data    json.RawMessage `json:"data"`
}

// Error includes additional metadata for a Supercharged error response.
type Error struct {
	Status  int
	Code    int
	Message string
	Data    interface{}
}

// Wrap wraps an error with default fields to conform to an Error struct.
func Wrap(err error) *Error {
	return &Error{
		Status:  400,
		Code:    CodeClientUndefined,
		Message: err.Error(),
		Data:    nil,
	}
}

// Error returns a string representation of the error.
func (e *Error) Error() string {
	return fmt.Sprintf("%s (%d)", e.Message, e.Code)
}
