package auth

import (
	"errors"
	"fmt"
	"time"
)

type FieldType int

const (
	UUIDFieldType FieldType = iota + 1
	TokenFieldType
	TimeRangeFieldType
	DurationFieldType
	ExpiredTokenFieldtype
	SecretKeyFieldType
	TokenCreationFieldType
)

var fieldNames = [7]string{
	"uuid",
	"token",
	"time range",
	"expired token",
	"secret key",
	"token creation",
}

func (f FieldType) String() string {
	if int(f) > len(fieldNames) {
		return "unknown"
	}
	return fieldNames[f-1]
}

type ValidationError struct {
	Field   FieldType
	Message string
}

func (e *ValidationError) Error() string {
	return fmt.Sprintf("{Field: %q, Message: %q}", e.Field, e.Message)
}

type ValidationErrors []error

func (ve ValidationErrors) Error() string {
	if len(ve) == 0 {
		return "validation errors:\n"
	}

	msg := "validation errors:\n"
	for _, err := range ve {
		msg += fmt.Sprintf("\t- %s\n", err)
	}
	return msg
}

func (ve *ValidationErrors) PushValidationError(err error) error {
	var validationErr *ValidationError
	if err == nil {
		return nil
	}
	if !errors.As(err, &validationErr) {
		return fmt.Errorf("unexpected error type: %T", err)
	}
	*ve = append(*ve, validationErr)
	return nil
}

func (ve ValidationErrors) Any() bool {
	return len(ve) > 0
}

func NewInvalidUUIDError(value string) error {
	return &ValidationError{
		Field:   UUIDFieldType,
		Message: fmt.Sprintf("invalid UUID: %s", value),
	}
}

func NewInvalidTokenError(reason string) error {
	return &ValidationError{
		Field:   TokenFieldType,
		Message: fmt.Sprintf("invalid token: %s", reason),
	}
}

func NewInvalidTimeRangeError(issuedAt, expiredAt time.Time) error {
	return &ValidationError{
		Field: TimeRangeFieldType,
		Message: fmt.Sprintf(
			"invalid time range: IssuedAt (%s) must be before ExpiredAt (%s)",
			issuedAt,
			expiredAt,
		),
	}
}

func NewInvalidDurationError(duration time.Duration, reason string) error {
	return &ValidationError{
		Field:   DurationFieldType,
		Message: fmt.Sprintf("invalid duration %v: %s", duration, reason),
	}
}

func NewExpiredTokenError(expiredAt time.Time) error {
	return &ValidationError{
		Field:   ExpiredTokenFieldtype,
		Message: fmt.Sprintf("invalid expired at time %s", expiredAt.String()),
	}
}

// Added NewTokenCreationError function
func NewTokenCreationError(reason string) error {
	return &ValidationError{
		Field:   TokenCreationFieldType,
		Message: fmt.Sprintf("failed to create token: %s", reason),
	}
}
