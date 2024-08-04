package outbound

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/teamkweku/code-odessey-hex-arch/internal/core/domain/user"
	mock_repo "github.com/teamkweku/code-odessey-hex-arch/internal/core/ports/outbound/mock"
	"go.uber.org/mock/gomock"
)

func TestUserRepository(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserRepo := mock_repo.NewMockUserRepository(ctrl)

	t.Run("GetUserByID", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		userID := uuid.New()
		expectedUser := &user.User{}

		mockUserRepo.EXPECT().GetUserByID(ctx, userID).
			Return(expectedUser, nil)

		resultUser, err := mockUserRepo.GetUserByID(ctx, userID)

		assert.NoError(t, err)
		assert.Equal(t, expectedUser, resultUser)
	})

	t.Run("GetUserByEmail", func(t *testing.T) {
		t.Parallel()
		ctx := context.Background()
		email := user.RandomEmailAddress(t)
		expectedUser := &user.User{}

		mockUserRepo.EXPECT().GetUserByEmail(ctx, email).
			Return(expectedUser, nil)

		resultUser, err := mockUserRepo.GetUserByEmail(ctx, email)
		assert.NoError(t, err)
		assert.Equal(t, expectedUser, resultUser)
	})

	t.Run("CreateUser", func(t *testing.T) {
		t.Parallel()
		ctx := context.Background()
		req := &user.RegistrationRequest{}
		expectedUser := &user.User{}

		mockUserRepo.EXPECT().CreateUser(ctx, req).
			Return(expectedUser, nil)

		resultUser, err := mockUserRepo.CreateUser(ctx, req)

		assert.NoError(t, err)
		assert.Equal(t, expectedUser, resultUser)
	})

	t.Run("UpdateUser", func(t *testing.T) {
		t.Parallel()
		ctx := context.Background()
		updateReq := &user.UpdateRequest{}
		expectedUser := &user.User{}

		mockUserRepo.EXPECT().UpdateUser(ctx, updateReq).
			Return(expectedUser, nil)

		resultUser, err := mockUserRepo.UpdateUser(ctx, updateReq)
		assert.NoError(t, err)
		assert.Equal(t, expectedUser, resultUser)
	})
}
