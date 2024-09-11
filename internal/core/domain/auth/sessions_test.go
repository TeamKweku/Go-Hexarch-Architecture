package auth

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewSessions(t *testing.T) {
	t.Parallel()

	userID := uuid.New()
	refreshToken := "refresh_token"
	userAgent := "test_user_agent"
	clientIP := "127.0.0.1"
	expiresAt := time.Now().Add(time.Hour)

	refreshPayload := &Payload{
		ExpiredAt: expiresAt,
	}

	session, err := NewSessions(userID, refreshToken, refreshPayload, userAgent, clientIP)

	require.NoError(t, err)
	assert.NotNil(t, session)

	assert.NotEqual(t, uuid.Nil, session.ID)
	assert.Equal(t, userID, session.UserID)
	assert.Equal(t, refreshToken, session.RefreshToken)
	assert.Equal(t, userAgent, session.UserAgent)
	assert.Equal(t, clientIP, session.ClientIP)
	assert.False(t, session.IsBlocked)
	assert.Equal(t, expiresAt, session.ExpiresAt)
	assert.WithinDuration(t, time.Now(), session.CreatedAt, time.Second)
}
