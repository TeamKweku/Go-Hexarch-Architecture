package inbound

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/teamkweku/code-odessey-hex-arch/internal/core/domain/user"
	mock_user "github.com/teamkweku/code-odessey-hex-arch/internal/core/ports/inbound/mock"
	"go.uber.org/mock/gomock"
)

func TestUserService_Register(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserService := mock_user.NewMockUserService(ctrl)

	t.Run("Register", func(t *testing.T) {
		t.Parallel()
		ctx := context.Background()
		req := &user.RegistrationRequest{}
		expectedUser := &user.User{}

		mockUserService.EXPECT().Register(ctx, req).
			Return(expectedUser, nil)

		resultUser, err := mockUserService.Register(ctx, req)

		assert.NoError(t, err)
		assert.Equal(t, expectedUser, resultUser)
	})

	t.Run("GetUser", func(t *testing.T) {
		t.Parallel()
		ctx := context.Background()
		userID := uuid.New()
		expectedUser := &user.User{}

		mockUserService.EXPECT().GetUser(ctx, userID).
			Return(expectedUser, nil)

		resultUser, err := mockUserService.GetUser(ctx, userID)

		assert.NoError(t, err)
		assert.Equal(t, expectedUser, resultUser)
	})

	t.Run("Authenticate", func(t *testing.T) {
		t.Parallel()
		ctx := context.Background()
		loginReq := &user.LoginRequest{}

		expectedUser := &user.User{}

		mockUserService.EXPECT().Authenticate(ctx, loginReq).
			Return(expectedUser, nil)

		resultUser, err := mockUserService.Authenticate(ctx, loginReq)

		assert.NoError(t, err)
		assert.Equal(t, expectedUser, resultUser)
	})

	t.Run("UpdateUser", func(t *testing.T) {
		t.Parallel()
		ctx := context.Background()
		updateReq := &user.UpdateRequest{}

		expectedUser := &user.User{}

		mockUserService.EXPECT().UpdateUser(ctx, updateReq).
			Return(expectedUser, nil)

		resultUser, err := mockUserService.UpdateUser(ctx, updateReq)

		assert.NoError(t, err)
		assert.Equal(t, expectedUser, resultUser)
	})
}
