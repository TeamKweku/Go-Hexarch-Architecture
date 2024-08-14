package user

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	domainUser "github.com/teamkweku/code-odessey-hex-arch/internal/core/domain/user"
	mock_outbound "github.com/teamkweku/code-odessey-hex-arch/internal/core/ports/outbound/mock"
)

var anyError = errors.New("any error")

func Test_service_Register(t *testing.T) {
	t.Parallel()

	req := domainUser.RandomRegistrationRequest(t)

	testCases := []struct {
		name     string
		wantUser *domainUser.User
		wantErr  error
	}{
		{
			name:     "repo call succeeds",
			wantUser: &domainUser.User{},
			wantErr:  nil,
		},
		{
			name:     "repo returns any error",
			wantUser: nil,
			wantErr:  anyError,
		},
	}

	for _, tc := range testCases {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockRepo := mock_outbound.NewMockUserRepository(ctrl)
			svc := NewUserService(mockRepo)

			mockRepo.EXPECT().
				CreateUser(gomock.Any(), req).
				Return(tc.wantUser, tc.wantErr)

			gotUser, gotErr := svc.Register(context.Background(), req)

			assert.Equal(t, tc.wantUser, gotUser)
			assert.ErrorIs(t, gotErr, tc.wantErr)
		})
	}
}

func Test_Service_Authenticate(t *testing.T) {
	t.Parallel()

	req := domainUser.RandomLoginRequest(t)

	testCases := []struct {
		name               string
		repoUser           *domainUser.User
		repoErr            error
		passwordComparator domainUser.PasswordComparator
		wantUser           *domainUser.User
		wantErr            error
	}{
		{
			name:     "success",
			repoUser: &domainUser.User{},
			repoErr:  nil,
			passwordComparator: func(hash domainUser.PasswordHash, candidate string) error {
				return nil
			},
			wantUser: &domainUser.User{},
			wantErr:  nil,
		},
		{
			name:               "repo returns NotFoundError",
			repoUser:           nil,
			repoErr:            &domainUser.NotFoundError{},
			passwordComparator: nil,
			wantUser:           nil,
			wantErr: &domainUser.AuthError{
				Cause: &domainUser.NotFoundError{},
			},
		},
		{
			name:               "repo returns any other error",
			repoUser:           nil,
			repoErr:            anyError,
			passwordComparator: nil,
			wantUser:           nil,
			wantErr:            anyError,
		},
		{
			name:     "passwordComparator returns an error",
			repoUser: &domainUser.User{},
			repoErr:  nil,
			passwordComparator: func(hash domainUser.PasswordHash, candidate string) error {
				return errors.New("invalid password")
			},
			wantUser: nil,
			wantErr: &domainUser.AuthError{
				Cause: errors.New("invalid password"),
			},
		},
	}

	for _, tc := range testCases {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockRepo := mock_outbound.NewMockUserRepository(ctrl)
			svc := &UserService{
				repo:               mockRepo,
				passwordComparator: tc.passwordComparator,
			}

			mockRepo.EXPECT().
				GetUserByEmail(gomock.Any(), req.Email()).
				Return(tc.repoUser, tc.repoErr)
			gotUser, gotErr := svc.Authenticate(context.Background(), req)

			assert.Equal(t, tc.wantUser, gotUser)
			assert.Equal(t, gotErr, tc.wantErr)
		})
	}
}

func Test_Service_GetUser(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name     string
		wantUser *domainUser.User
		wantErr  error
	}{
		{
			name:     "repo call succeeds",
			wantUser: &domainUser.User{},
			wantErr:  nil,
		},
		{
			name:     "repo returns any error",
			wantUser: nil,
			wantErr:  anyError,
		},
	}

	for _, tc := range testCases {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockRepo := mock_outbound.NewMockUserRepository(ctrl)
			svc := NewUserService(mockRepo)
			userID := uuid.New()

			mockRepo.EXPECT().
				GetUserByID(gomock.Any(), userID).
				Return(tc.wantUser, tc.wantErr)

			gotUser, gotErr := svc.GetUser(context.Background(), userID)

			assert.Equal(t, tc.wantUser, gotUser)
			assert.ErrorIs(t, gotErr, tc.wantErr)
		})
	}
}

func Test_service_UpdateUser(t *testing.T) {
	t.Parallel()

	req := domainUser.RandomUpdateRequest(t)

	testCases := []struct {
		name     string
		wantUser *domainUser.User
		wantErr  error
	}{
		{
			name:     "repo call succeeds",
			wantUser: &domainUser.User{},
			wantErr:  nil,
		},
		{
			name:     "repo returns any error",
			wantUser: nil,
			wantErr:  anyError,
		},
	}

	for _, tc := range testCases {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockRepo := mock_outbound.NewMockUserRepository(ctrl)
			svc := NewUserService(mockRepo)

			mockRepo.EXPECT().
				UpdateUser(gomock.Any(), req).
				Return(tc.wantUser, tc.wantErr)

			gotUser, gotErr := svc.UpdateUser(context.Background(), req)

			assert.Equal(t, tc.wantUser, gotUser)
			if tc.wantErr != nil {
				assert.Error(t, gotErr)
				assert.True(
					t,
					errors.Is(gotErr, tc.wantErr),
					"expected error %v, got %v",
					tc.wantErr,
					gotErr,
				)
			} else {
				assert.NoError(t, gotErr)
			}
		})
	}
}
