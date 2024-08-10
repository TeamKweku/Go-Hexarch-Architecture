package sqlc

import (
	"context"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/stretchr/testify/require"

	"github.com/teamkweku/code-odessey-hex-arch/internal/core/domain/user"
)

//nolint:paralleltest
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
	require.NotZero(t, user.ID)
	require.Equal(t, user.PasswordChangedAt.UTC(), time.Time{}.UTC())
	require.Equal(t, user.UpdatedAt.UTC(), time.Time{}.UTC())
}

//nolint:paralleltest
func TestGetUserByEmail(t *testing.T) {
	// create a user
	randomUser := user.RandomUser(t)

	arg := CreateUserParams{
		Username:     randomUser.Username().String(),
		Etag:         randomUser.ETag().String(),
		Email:        randomUser.Email().String(),
		PasswordHash: string(randomUser.PasswordHash().Bytes()),
	}

	user, err := testQueries.CreateUser(context.Background(), arg)
	require.NoError(t, err)

	// get user by email
	fetchedUser, err := testQueries.GetUserByEmail(
		context.Background(),
		user.Email,
	)
	require.NoError(t, err)
	require.NotEmpty(t, fetchedUser)

	require.Equal(t, user.ID, fetchedUser.ID)
	require.Equal(t, user.Username, fetchedUser.Username)
	require.Equal(t, user.PasswordHash, fetchedUser.PasswordHash)
	require.True(t, fetchedUser.UpdatedAt.IsZero())

	require.WithinDuration(
		t,
		user.CreatedAt,
		fetchedUser.CreatedAt,
		time.Second,
	)
}

// test getting user by userID
//
//nolint:paralleltest
func TestGetUserByID(t *testing.T) {
	// create a user
	randomUser := user.RandomUser(t)

	arg := CreateUserParams{
		Username:     randomUser.Username().String(),
		Etag:         randomUser.ETag().String(),
		Email:        randomUser.Email().String(),
		PasswordHash: string(randomUser.PasswordHash().Bytes()),
	}

	user, err := testQueries.CreateUser(context.Background(), arg)
	require.NoError(t, err)

	// get user by ID
	fetchedUser, err := testQueries.GetUserById(context.Background(), user.ID)
	require.NoError(t, err)
	require.NotEmpty(t, fetchedUser)

	require.Equal(t, user.Email, fetchedUser.Email)
	require.Equal(t, user.Username, fetchedUser.Username)
	require.Equal(t, user.PasswordHash, fetchedUser.PasswordHash)
	require.True(t, fetchedUser.UpdatedAt.IsZero())

	require.WithinDuration(
		t,
		user.CreatedAt,
		fetchedUser.CreatedAt,
		time.Second,
	)
}

// testing update user function with all parameters
//
//nolint:paralleltest
func TestUpdateUser(t *testing.T) {
	// create a user
	randomUser := user.RandomUser(t)
	arg := CreateUserParams{
		Username:     randomUser.Username().String(),
		Etag:         randomUser.ETag().String(),
		Email:        randomUser.Email().String(),
		PasswordHash: string(randomUser.PasswordHash().Bytes()),
	}
	newUser, err := testQueries.CreateUser(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, newUser)

	// create update parameters
	updateUser := user.RandomUser(t)
	updateUserName := updateUser.Username().String()
	updateUserEmail := updateUser.Email().String()
	updateUserPasswordHash := string(updateUser.PasswordHash().Bytes())
	updateArg := UpdateUserParams{
		ID: newUser.ID, // Use the ID of the user we just created
		Username: pgtype.Text{
			String: updateUserName,
			Valid:  true,
		},
		Email: pgtype.Text{
			String: updateUserEmail,
			Valid:  true,
		},
		PasswordHash: pgtype.Text{
			String: updateUserPasswordHash,
			Valid:  true,
		},
	}
	updatedUser, err := testQueries.UpdateUser(context.Background(), updateArg)
	require.NoError(t, err)
	require.NotEmpty(t, updatedUser)

	// Add more assertions to verify the update was successful
	require.Equal(t, newUser.ID, updatedUser.ID)
	require.Equal(t, updateUserName, updatedUser.Username)
	require.Equal(t, updateUserEmail, updatedUser.Email)
	require.Equal(t, updateUserPasswordHash, updatedUser.PasswordHash)
	require.WithinDuration(
		t,
		time.Now().UTC(),
		updatedUser.UpdatedAt.UTC(),
		2*time.Second,
	)
}

//nolint:paralleltest
func TestUpdateUserPartial(t *testing.T) {
	randomUser := user.RandomUser(t)
	arg := CreateUserParams{
		Username:     randomUser.Username().String(),
		Etag:         randomUser.ETag().String(),
		Email:        randomUser.Email().String(),
		PasswordHash: string(randomUser.PasswordHash().Bytes()),
	}
	newUser, err := testQueries.CreateUser(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, newUser)

	// new email
	newEmail := user.RandomEmailAddress(t).String()

	// create update parameters
	updateArg := UpdateUserParams{
		ID: newUser.ID, // Use the ID of the user we just created
		Email: pgtype.Text{
			String: newEmail,
			Valid:  true,
		},
	}
	updatedUser, err := testQueries.UpdateUser(context.Background(), updateArg)
	require.NoError(t, err)
	require.NotEmpty(t, updatedUser)

	// Add more assertions to verify the update was successful
	require.Equal(t, newUser.ID, updatedUser.ID)
	require.Equal(t, newUser.Username, updatedUser.Username)
	require.NotEqual(t, newUser.Email, updatedUser.Email)
	require.Equal(t, newUser.Role, updatedUser.Role)
	require.Equal(t, newUser.Etag, updatedUser.Etag)
	require.Equal(t, newUser.PasswordHash, updatedUser.PasswordHash)
	require.Equal(t, newUser.CreatedAt, updatedUser.CreatedAt)
	require.WithinDuration(
		t,
		time.Now().UTC(),
		updatedUser.UpdatedAt.UTC(),
		2*time.Second,
	)
}
