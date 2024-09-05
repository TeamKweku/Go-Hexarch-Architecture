package auth

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/teamkweku/code-odessey-hex-arch/config"
	"github.com/teamkweku/code-odessey-hex-arch/internal/core/domain/auth"
	"github.com/teamkweku/code-odessey-hex-arch/internal/core/domain/user"
)

func TestPasetoToken_CreateToken(t *testing.T) {
	t.Parallel()

	pt, err := NewPasetoToken()
	require.NoError(t, err)

	usr := user.RandomUser(t)
	cfg := config.Config{
		AccessTokenDuration: "30m",
	}

	token, payload, err := pt.CreateToken(usr, cfg)
	require.NoError(t, err)
	assert.NotEmpty(t, token)
	assert.NotNil(t, payload)
}

func TestPasetoToken_VerifyToken(t *testing.T) {
	t.Parallel()

	pt, err := NewPasetoToken()
	require.NoError(t, err)

	usr := user.RandomUser(t)
	cfg := config.Config{
		AccessTokenDuration: "30m",
	}

	token, payload, err := pt.CreateToken(usr, cfg)
	require.NoError(t, err)
	assert.NotEmpty(t, token)
	assert.NotNil(t, payload)

	verifiedPayload, err := pt.VerifyToken(token)
	require.NoError(t, err)
	assert.Equal(t, payload, verifiedPayload)
}

func TestPasetoToken_VerifyToken_Expired(t *testing.T) {
	t.Parallel()

	pt, err := NewPasetoToken()
	require.NoError(t, err)

	usr := user.RandomUser(t)
	cfg := config.Config{
		AccessTokenDuration: "10ms", // Set a short duration for testing expiration
	}

	token, payload, err := pt.CreateToken(usr, cfg)
	require.NoError(t, err)
	assert.NotEmpty(t, token)
	assert.NotNil(t, payload)

	// Simulate token expiration
	time.Sleep(15 * time.Millisecond)

	_, err = pt.VerifyToken(token)
	assert.Error(t, err)
	var expiredErr *auth.ValidationError
	assert.ErrorAs(t, err, &expiredErr)
	assert.Equal(t, auth.ExpiredTokenFieldtype, expiredErr.Field)
}
