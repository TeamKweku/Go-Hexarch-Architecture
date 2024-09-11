package auth

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/teamkweku/code-odessey-hex-arch/internal/core/domain/user"
)

func TestNewPayload(t *testing.T) {
	t.Parallel()

	testUser := user.RandomUser(t)
	durationStr := "1h"

	payload, err := NewPayload(testUser, durationStr)
	require.NoError(t, err)

	assert.NotNil(t, payload)
	assert.NotEqual(t, uuid.Nil, payload.ID)
	assert.Equal(t, testUser.ID(), payload.UserID)
	assert.Equal(t, testUser.Role().String(), payload.Role)
	assert.WithinDuration(t, time.Now(), payload.IssuedAt, time.Second)
	assert.WithinDuration(t, time.Now().Add(time.Hour), payload.ExpiredAt, time.Second)
}

func TestPayload_Valid(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name        string
		expiredAt   time.Time
		expectedErr error
	}{
		{
			name:        "Valid payload",
			expiredAt:   time.Now().Add(time.Hour),
			expectedErr: nil,
		},
		{
			name:        "Expired payload",
			expiredAt:   time.Now().Add(-time.Hour),
			expectedErr: &ValidationError{Field: ExpiredTokenFieldtype},
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			payload := &Payload{
				ID:        uuid.New(),
				UserID:    uuid.New(),
				Role:      "user",
				IssuedAt:  time.Now(),
				ExpiredAt: tc.expiredAt,
			}

			err := payload.Valid()
			if tc.expectedErr == nil {
				assert.NoError(t, err)
			} else {
				assert.IsType(t, tc.expectedErr, err)
			}
		})
	}
}

func TestNewPayload_InvalidDuration(t *testing.T) {
	t.Parallel()

	testUser := user.RandomUser(t)
	durationStr := "invalid"

	payload, err := NewPayload(testUser, durationStr)
	assert.Error(t, err)
	assert.Nil(t, payload)
	assert.IsType(t, &ValidationError{}, err)
	validationErr := err.(*ValidationError)
	assert.Equal(t, DurationFieldType, validationErr.Field)
}
