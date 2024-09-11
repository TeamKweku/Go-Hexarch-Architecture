package auth

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFieldType_String(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name      string
		fieldType FieldType
		expected  string
	}{
		{"UUIDFieldType", UUIDFieldType, "uuid"},
		{"TokenFieldType", TokenFieldType, "token"},
		{"TimeRangeFieldType", TimeRangeFieldType, "time range"},
		{"DurationFieldType", DurationFieldType, "duration"},
		{"ExpiredTokenFieldtype", ExpiredTokenFieldtype, "expired token"},
		{"SecretKeyFieldType", SecretKeyFieldType, "secret key"},
		{"TokenCreationFieldType", TokenCreationFieldType, "token creation"},
		{"Unknown FieldType", FieldType(100), "unknown"},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tc.expected, tc.fieldType.String())
		})
	}
}

func TestValidationError_Error(t *testing.T) {
	t.Parallel()

	err := &ValidationError{
		Field:   UUIDFieldType,
		Message: "invalid UUID",
	}

	assert.Equal(t, `{Field: "uuid", Message: "invalid UUID"}`, err.Error())
}

func TestValidationErrors_Error(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name     string
		errors   ValidationErrors
		expected string
	}{
		{
			name:     "Empty errors",
			errors:   ValidationErrors{},
			expected: "validation errors:\n",
		},
		{
			name: "Single error",
			errors: ValidationErrors{
				&ValidationError{Field: UUIDFieldType, Message: "invalid UUID"},
			},
			expected: "validation errors:\n\t- {Field: \"uuid\", Message: \"invalid UUID\"}\n",
		},
		{
			name: "Multiple errors",
			errors: ValidationErrors{
				&ValidationError{Field: UUIDFieldType, Message: "invalid UUID"},
				&ValidationError{Field: TokenFieldType, Message: "invalid token"},
			},
			expected: "validation errors:\n\t- {Field: \"uuid\", Message: \"invalid UUID\"}\n\t- {Field: \"token\", Message: \"invalid token\"}\n",
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tc.expected, tc.errors.Error())
		})
	}
}

func TestValidationErrors_PushValidationError(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name          string
		initialErrors ValidationErrors
		newError      error
		expectedLen   int
		expectedError error
	}{
		{
			name:          "Add valid ValidationError",
			initialErrors: ValidationErrors{},
			newError:      &ValidationError{Field: UUIDFieldType, Message: "invalid UUID"},
			expectedLen:   1,
			expectedError: nil,
		},
		{
			name:          "Add nil error",
			initialErrors: ValidationErrors{},
			newError:      nil,
			expectedLen:   0,
			expectedError: nil,
		},
		{
			name:          "Add non-ValidationError",
			initialErrors: ValidationErrors{},
			newError:      errors.New("regular error"),
			expectedLen:   0,
			expectedError: errors.New("unexpected error type: *errors.errorString"),
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			ve := tc.initialErrors
			err := ve.PushValidationError(tc.newError)
			assert.Equal(t, tc.expectedLen, len(ve))
			if tc.expectedError == nil {
				assert.NoError(t, err)
			} else {
				assert.EqualError(t, err, tc.expectedError.Error())
			}
		})
	}
}

func TestNewTokenCreationError(t *testing.T) {
	t.Parallel()

	err := NewTokenCreationError("test reason")
	validationErr, ok := err.(*ValidationError)
	assert.True(t, ok)
	assert.Equal(t, TokenCreationFieldType, validationErr.Field)
	assert.Equal(t, "failed to create token: test reason", validationErr.Message)
}
