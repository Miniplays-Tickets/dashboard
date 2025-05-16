package validator

import "errors"

var (
	ErrValidationFailed    = errors.New("Validierung fehlgeschlagen")
	ErrMaximumSizeExceeded = errors.New("Maximale Größe überschritten")
)
