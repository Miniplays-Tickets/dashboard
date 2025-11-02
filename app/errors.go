package app

import (
	"errors"
	"fmt"

	"github.com/TicketsBot-cloud/gdl/rest/request"
)

type ApiError struct {
	InternalError   error
	ExternalMessage string
}

var _ error = (*ApiError)(nil)

func NewError(internalError error, externalMessage string) *ApiError {
	return &ApiError{
		InternalError:   internalError,
		ExternalMessage: externalMessage,
	}
}

func NewServerError(internalError error) *ApiError {
	var restError request.RestError
	if errors.As(internalError, &restError) {
		return NewError(internalError, restError.Error())
	}

	return NewError(internalError, "Ein interner Server Fehler ist aufgetreten!")
}

func (e *ApiError) Error() string {
	var restError request.RestError
	if errors.As(e.InternalError, &restError) {
		return fmt.Sprintf("Interner Fehler: %v, Externe Nachricht: %s, rest error: Discord returned HTTP %d: %s",
			e.InternalError, e.ExternalMessage, restError.StatusCode, restError.ApiError.Message)
	} else {
		return fmt.Sprintf("Interner Fehler: %v, Externe Nachricht: %s", e.InternalError, e.ExternalMessage)
	}
}

func (e *ApiError) Unwrap() error {
	return e.InternalError
}
