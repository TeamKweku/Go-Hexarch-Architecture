package sqlc

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/teamkweku/code-odessey-hex-arch/internal/core/domain/user"
)

func TestCreateUser(t *testing.T) {
	randomUser := user.RandomUser(t)

	arg := CreateUserParams{
		Username:     randomUser.Username().String(),
		Email:        randomUser.Email().String(),
		PasswordHash: string(randomUser.PasswordHash().Bytes()),
		Etag:         randomUser.ETag().String(),
	}

	user, err := testQueries.CreateUser(context.Background(), arg)

	require.NoError(t, err)
	require.NotEmpty(t, user)

	require.Equal(t, arg.Username, user.Username)
	require.Equal(t, arg.Email, user.Email)
	require.Equal(t, arg.PasswordHash, user.PasswordHash)
	require.NotEmpty(t, user.Role)

}
