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
	}

	user, err := testQueries.CreateUser(context.Background(), arg)

	require.NoError(t, err)
	require.NotEmpty(t, user)
}
