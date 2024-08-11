package user

import (
	"errors"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/teamkweku/code-odessey-hex-arch/pkg/etag"
	"github.com/teamkweku/code-odessey-hex-arch/pkg/option"
)

func Test_ParseEmailAddress(t *testing.T) {
	t.Parallel()

	validEmailCandidate := RandomEmailAddressCandidate()

	testCases := []struct {
		name             string
		candidate        string
		wantEmailAddress EmailAddress
		wantErr          error
	}{
		{
			name:             "valid email address",
			candidate:        validEmailCandidate,
			wantEmailAddress: EmailAddress{raw: validEmailCandidate},
			wantErr:          nil,
		},
		{
			name:             "invalid email address",
			candidate:        "test",
			wantEmailAddress: EmailAddress{},
			wantErr:          NewEmailAddressFormatError("test"),
		},
	}

	for _, tc := range testCases {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			gotEmailAddress, gotErr := ParseEmailAddress(tc.candidate)

			assert.Equal(t, tc.wantEmailAddress, gotEmailAddress)
			assert.Equal(t, tc.wantErr, gotErr)
		})
	}
}

func Test_EmailAddress_String(t *testing.T) {
	t.Parallel()

	email := RandomEmailAddress(t)

	assert.Equal(t, email.raw, email.String())
}

func Test_ParseUsername(t *testing.T) {
	t.Parallel()

	validUsernameCandidate := RandomUsernameCandidate()

	testCases := []struct {
		name         string
		candidate    string
		wantUsername Username
		wantErr      error
	}{
		{
			name:         "valid username",
			candidate:    validUsernameCandidate,
			wantUsername: Username{raw: validUsernameCandidate},
			wantErr:      nil,
		},
		{
			name:         "username too short",
			candidate:    strings.Repeat("a", UsernameMinLen-1),
			wantUsername: Username{},
			wantErr:      NewUsernameTooShortError(),
		},
		{
			name:         "username too long",
			candidate:    strings.Repeat("a", UsernameMaxLen+1),
			wantUsername: Username{},
			wantErr:      NewUsernameTooLongError(),
		},
		{
			name:         "username contains invalid characters",
			candidate:    "test!",
			wantUsername: Username{},
			wantErr:      NewUsernameFormatError(),
		},
	}

	for _, tc := range testCases {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			gotUsername, gotErr := ParseUsername(tc.candidate)

			assert.Equal(t, tc.wantUsername, gotUsername)
			assert.Equal(t, tc.wantErr, gotErr)
		})
	}
}

func Test_Username_String(t *testing.T) {
	t.Parallel()

	username := RandomUsername(t)

	assert.Equal(t, username.raw, username.String())
}

func Test_parsePassword(t *testing.T) {
	t.Parallel()

	anyError := errors.New("any error")

	assertEmptyPasswordHash := func(t *testing.T, hash PasswordHash, candidate string) {
		t.Helper()
		assert.Empty(t, hash)
	}

	testCases := []struct {
		name               string
		candidate          string
		hasher             passwordHasher
		assertPasswordHash func(t *testing.T, hash PasswordHash, candidate string)
		wantErr            error
	}{
		{
			name:      "valid password",
			candidate: RandomPasswordCandidate(),
			hasher:    bcryptHash,
			assertPasswordHash: func(t *testing.T, hash PasswordHash, candidate string) {
				t.Helper()
				assert.NoError(t, BcryptCompare(hash, candidate))
			},
			wantErr: nil,
		},
		{
			name:               "password too short",
			candidate:          strings.Repeat("a", PasswordMinLen-1),
			hasher:             bcryptHash,
			assertPasswordHash: assertEmptyPasswordHash,
			wantErr:            NewPasswordTooShortError(),
		},
		{
			name:               "password too long",
			candidate:          strings.Repeat("a", PasswordMaxLen+1),
			hasher:             bcryptHash,
			assertPasswordHash: assertEmptyPasswordHash,
			wantErr:            NewPasswordTooLongError(),
		},
		{
			name:      "hasher returns any error",
			candidate: RandomPasswordCandidate(),
			hasher: func(secret string) (PasswordHash, error) {
				return PasswordHash{}, anyError
			},
			assertPasswordHash: assertEmptyPasswordHash,
			wantErr:            anyError,
		},
	}

	for _, tc := range testCases {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			gotPasswordHash, gotErr := parsePassword(tc.candidate, tc.hasher)

			tc.assertPasswordHash(t, gotPasswordHash, tc.candidate)
			assert.ErrorIs(t, gotErr, tc.wantErr)
		})
	}
}

func Test_NewPasswordHashFromTrustedSource(t *testing.T) {
	t.Parallel()

	want := RandomPasswordHash(t)

	got := NewPasswordHashFromTrustedSource(want.bytes)

	assert.Equal(t, want, got)
}

func Test_PasswordHash_GoString(t *testing.T) {
	t.Parallel()

	hash := RandomPasswordHash(t)

	assert.Equal(t, "PasswordHash{bytes:REDACTED}", hash.GoString())
}

func Test_PasswordHash_String(t *testing.T) {
	t.Parallel()

	hash := RandomPasswordHash(t)

	assert.Equal(t, "{REDACTED}", hash.String())
}

func TestParseRole(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name         string
		candidate    int
		expectedRole Role
		expectError  bool
	}{
		{"Valid Reader", int(RoleReader), RoleReader, false},
		{"Valid Author", int(RoleAuthor), RoleAuthor, false},
		{"Valid Editor", int(RoleEditor), RoleEditor, false},
		{"Valid Admin", int(RoleAdmin), RoleAdmin, false},
		{"Invalid Low", -1, Role(0), true},
		{"Invalid High", int(RoleAdmin) + 1, Role(0), true},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			role, err := ParseRole(tt.candidate)

			if tt.expectError {
				assert.Error(t, err)
				assert.Equal(t, Role(0), role)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedRole, role)
			}
		})
	}
}

func TestRole_String(t *testing.T) {
	t.Parallel()
	tests := []struct {
		role     Role
		expected string
	}{
		{RoleReader, "Reader"},
		{RoleAuthor, "Author"},
		{RoleEditor, "Editor"},
		{RoleAdmin, "Admin"},
		{Role(99), "Unknown"},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.expected, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tt.expected, tt.role.String())
		})
	}
}

func TestRole_GoString(t *testing.T) {
	t.Parallel()
	tests := []struct {
		role     Role
		expected string
	}{
		{RoleReader, "Role(0)"},
		{RoleAuthor, "Role(1)"},
		{RoleEditor, "Role(2)"},
		{RoleAdmin, "Role(3)"},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(fmt.Sprintf("Role(%d)", int(tt.role)), func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tt.expected, tt.role.GoString())
		})
	}
}

func TestRoleConstants(t *testing.T) {
	t.Parallel()
	assert.Equal(t, Role(0), RoleReader)
	assert.Equal(t, Role(1), RoleAuthor)
	assert.Equal(t, Role(2), RoleEditor)
	assert.Equal(t, Role(3), RoleAdmin)
}

func Test_NewUser(t *testing.T) {
	t.Parallel()

	id := uuid.New()
	eTag := etag.Random()
	username := RandomUsername(t)
	email := RandomEmailAddress(t)
	passwordHash := RandomPasswordHash(t)
	role := RandomRole(t)
	createdAt := time.Now()
	passwordChangedAt := createdAt.Add(time.Hour)

	user := NewUser(
		id,
		eTag,
		username,
		email,
		passwordHash,
		role,
		createdAt,
		passwordChangedAt,
	)

	assert.NotNil(t, user)
	assert.Equal(t, id, user.ID())
	assert.Equal(t, eTag, user.ETag())
	assert.Equal(t, username, user.Username())
	assert.Equal(t, email, user.Email())
	assert.Equal(t, passwordHash, user.PasswordHash())
}

func Test_User_GoString(t *testing.T) {
	t.Parallel()

	user := RandomUser(t)

	expected := fmt.Sprintf(
		"User{id:%#v, eTag:%#v, username:%#v, email:%#v, passwordHash:%#v,",
		user.id, user.eTag, user.username, user.email, user.passwordHash,
	)

	assert.Equal(t, expected, user.GoString())
}

func Test_User_String(t *testing.T) {
	t.Parallel()

	user := RandomUser(t)

	expected := fmt.Sprintf(
		"{ %s %s %s %s %s %s}",
		user.id,
		user.eTag,
		user.username,
		user.email,
		user.passwordHash,
		user.role,
	)

	assert.Equal(t, expected, user.String())
}

func Test_NewRegistrationRequest(t *testing.T) {
	t.Parallel()

	username := RandomUsername(t)
	email := RandomEmailAddress(t)
	passwordHash := RandomPasswordHash(t)

	req := NewRegistrationRequest(username, email, passwordHash)

	assert.NotNil(t, req)
	assert.Equal(t, username, req.Username())
	assert.Equal(t, email, req.Email())
	assert.Equal(t, passwordHash, req.PasswordHash())
}

func Test_ParseRegistrationRequest(t *testing.T) {
	t.Parallel()

	validUsername := RandomUsernameCandidate()
	validEmail := RandomEmailAddressCandidate()
	validPassword := RandomPasswordCandidate()

	testCases := []struct {
		name      string
		username  string
		email     string
		password  string
		wantError bool
	}{
		{
			name:      "valid request",
			username:  validUsername,
			email:     validEmail,
			password:  validPassword,
			wantError: false,
		},
		{
			name:      "invalid username",
			username:  "a",
			email:     validEmail,
			password:  validPassword,
			wantError: true,
		},
		{
			name:      "invalid email",
			username:  validUsername,
			email:     "invalid-email",
			password:  validPassword,
			wantError: true,
		},
		{
			name:      "invalid password",
			username:  validUsername,
			email:     validEmail,
			password:  "short",
			wantError: true,
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			req, err := ParseRegistrationRequest(
				tc.username,
				tc.email,
				tc.password,
			)

			if tc.wantError {
				assert.Error(t, err)
				assert.Nil(t, req)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, req)
				assert.Equal(t, tc.username, req.Username().String())
				assert.Equal(t, tc.email, req.Email().String())
				assert.NoError(t, BcryptCompare(req.PasswordHash(), tc.password))
			}
		})
	}
}

func Test_RegistrationRequest_Equal(t *testing.T) {
	t.Parallel()

	// Create a password and hash it
	password := RandomPasswordCandidate()
	passwordHash, err := ParsePassword(password)
	require.NoError(t, err)

	// Create req1 with the hashed password
	username := RandomUsername(t)
	email := RandomEmailAddress(t)
	req1 := NewRegistrationRequest(username, email, passwordHash)

	// Create req2 with different data
	req2 := RandomRegistrationRequest(t)

	// Create req3 as a copy of req1
	req3 := NewRegistrationRequest(
		req1.Username(),
		req1.Email(),
		req1.PasswordHash(),
	)

	testCases := []struct {
		name     string
		req1     *RegistrationRequest
		req2     *RegistrationRequest
		password string
		want     bool
	}{
		{
			name:     "equal requests",
			req1:     req1,
			req2:     req3,
			password: password,
			want:     true,
		},
		{
			name:     "different requests",
			req1:     req1,
			req2:     req2,
			password: password,
			want:     false,
		},
		{
			name:     "wrong password",
			req1:     req1,
			req2:     req3,
			password: "wrong_password",
			want:     false,
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			got := tc.req1.Equal(tc.req2, tc.password)
			assert.Equal(t, tc.want, got)
		})
	}
}

func Test_ParseLoginRequest(t *testing.T) {
	t.Parallel()

	validEmail := RandomEmailAddressCandidate()
	validPassword := RandomPasswordCandidate()

	testCases := []struct {
		name      string
		email     string
		password  string
		wantError bool
	}{
		{
			name:      "valid request",
			email:     validEmail,
			password:  validPassword,
			wantError: false,
		},
		{
			name:      "invalid email",
			email:     "invalid-email",
			password:  validPassword,
			wantError: true,
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			req, err := ParseLoginRequest(tc.email, tc.password)

			if tc.wantError {
				assert.Error(t, err)
				assert.Nil(t, req)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, req)
				assert.Equal(t, tc.email, req.Email().String())
				assert.Equal(t, tc.password, req.PasswordCandidate())
			}
		})
	}
}

func Test_LoginRequest_GoString(t *testing.T) {
	t.Parallel()

	req := RandomLoginRequest(t)

	expected := fmt.Sprintf(
		"LoginRequest{email:%#v, passwordCandidate:REDACTED}",
		req.Email(),
	)

	assert.Equal(t, expected, req.GoString())
}

func Test_LoginRequest_String(t *testing.T) {
	t.Parallel()

	req := RandomLoginRequest(t)

	expected := fmt.Sprintf("{%s REDACTED}", req.Email())

	assert.Equal(t, expected, req.String())
}

func Test_NewUpdateRequest(t *testing.T) {
	t.Parallel()

	userID := uuid.New()
	eTag := etag.Random()
	email := option.Some(RandomEmailAddress(t))
	passwordHash := option.Some(RandomPasswordHash(t))
	username := option.Some(RandomUsername(t))
	role := option.Some(RandomRole(t))

	req := NewUpdateRequest(userID, eTag, email, passwordHash, username, role)

	assert.NotNil(t, req)
	assert.Equal(t, userID, req.UserID())
	assert.Equal(t, eTag, req.ETag())
	assert.Equal(t, email, req.Email())
	assert.Equal(t, passwordHash, req.PasswordHash())
	assert.Equal(t, username, req.Username())
	assert.Equal(t, role, req.Role())
}

// func Test_ParseUpdateRequest(t *testing.T) {
// 	t.Parallel()

// 	userID := uuid.New()
// 	eTag := etag.Random()
// 	validEmail := RandomEmailAddressCandidate()
// 	validPassword := RandomPasswordCandidate()
// 	validUsername := RandomUsernameCandidate()
// 	validRole := int(RandomRole(t))
// 	email, err := ParseEmailAddress(validEmail)
// 	require.NoError(t, err)
// 	passwordHash, err := ParsePassword(validPassword)
// 	require.NoError(t, err)

// 	assertEqualUpdateRequest := func(t *testing.T, want, got *UpdateRequest) {
// 		t.Helper()

// 		if want == nil && got == nil {
// 			return
// 		}

// 		assert.Equal(t, want.userID, got.userID)
// 		assert.Equal(t, want.email, got.email)
// 		assert.Equal(t, want.username, got.username)
// 		assert.Equal(t, want.role, got.role)

// 		if !want.passwordHash.IsSome() {
// 			assert.True(t, !got.passwordHash.IsSome(),
// 				"passwordHash should be an empty Option, but value %v was found",
// 				got.passwordHash.UnwrapOrZero())
// 		} else {
// 			err := BcryptCompare(got.passwordHash.UnwrapOrZero(), validPassword)
// 			assert.NoError(t, err)
// 		}
// 	}

// 	testCases := []struct {
// 		name               string
// 		emailCandidate     option.Option[string]
// 		passwordCandidate  option.Option[string]
// 		username           option.Option[string]
// 		role               option.Option[string]
// 		wantedUpdteRequest *UpdateRequest
// 		wantErr            error
// 	}{
// 		{
// 			name:              "valid inputs, optional inputs present",
// 			emailCandidate:    option.Some(validEmail),
// 			passwordCandidate: option.Some(validPassword),
// 			username:          option.Some(validUsername),
// 			role:              option.Some(string(validRole)),
// 			wantedUpdteRequest: &UpdateRequest{
// 				userID:       userID,
// 				eTag:         eTag,
// 				username: option.Some(),
// 			},
// 		},
// 	}

// }
func Test_ParseUpdateRequest(t *testing.T) {
	t.Parallel()

	validUserID := uuid.New()
	validETag := etag.Random()
	validEmailCandidate := option.Some(RandomEmailAddressCandidate())
	validPasswordCandidate := option.Some(RandomPasswordCandidate())
	validUsernameCandidate := option.Some(RandomUsernameCandidate())
	validRoleCandidate := option.Some(int(RandomRole(t)))

	testCases := []struct {
		name              string
		userID            uuid.UUID
		eTag              etag.ETag
		emailCandidate    option.Option[string]
		passwordCandidate option.Option[string]
		usernameCandidate option.Option[string]
		roleCandidate     option.Option[int]
		wantErr           bool
	}{
		{
			name:              "valid update request",
			userID:            validUserID,
			eTag:              validETag,
			emailCandidate:    validEmailCandidate,
			passwordCandidate: validPasswordCandidate,
			usernameCandidate: validUsernameCandidate,
			roleCandidate:     validRoleCandidate,
			wantErr:           false,
		},
		{
			name:              "invalid email",
			userID:            validUserID,
			eTag:              validETag,
			emailCandidate:    option.Some("invalid-email"),
			passwordCandidate: validPasswordCandidate,
			usernameCandidate: validUsernameCandidate,
			roleCandidate:     validRoleCandidate,
			wantErr:           true,
		},
		{
			name:              "invalid password",
			userID:            validUserID,
			eTag:              validETag,
			emailCandidate:    validEmailCandidate,
			passwordCandidate: option.Some("short"),
			usernameCandidate: validUsernameCandidate,
			roleCandidate:     validRoleCandidate,
			wantErr:           true,
		},
		{
			name:              "invalid username",
			userID:            validUserID,
			eTag:              validETag,
			emailCandidate:    validEmailCandidate,
			passwordCandidate: validPasswordCandidate,
			usernameCandidate: option.Some("a"),
			roleCandidate:     validRoleCandidate,
			wantErr:           true,
		},
		{
			name:              "invalid role",
			userID:            validUserID,
			eTag:              validETag,
			emailCandidate:    validEmailCandidate,
			passwordCandidate: validPasswordCandidate,
			usernameCandidate: validUsernameCandidate,
			roleCandidate:     option.Some(-1),
			wantErr:           true,
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			req, err := ParseUpdateRequest(
				tc.userID,
				tc.eTag,
				tc.emailCandidate,
				tc.passwordCandidate,
				tc.usernameCandidate,
				tc.roleCandidate,
			)

			if tc.wantErr {
				assert.Error(t, err)
				assert.Nil(t, req)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, req)
				assert.Equal(t, tc.userID, req.UserID())
				assert.Equal(t, tc.eTag, req.ETag())
				assert.Equal(t, tc.emailCandidate, option.Some(req.Email().UnwrapOrZero().String()))
				assert.Equal(t, tc.usernameCandidate, option.Some(req.Username().UnwrapOrZero().String()))
				assert.Equal(t, tc.roleCandidate, option.Some(int(req.Role().UnwrapOrZero())))
			}
		})
	}
}

func Test_UpdateRequest_UserID(t *testing.T) {
	t.Parallel()

	userID := uuid.New()
	req := &UpdateRequest{userID: userID}

	assert.Equal(t, userID, req.UserID())
}

func Test_UpdateRequest_ETag(t *testing.T) {
	t.Parallel()

	eTag := etag.Random()
	req := &UpdateRequest{eTag: eTag}

	assert.Equal(t, eTag, req.ETag())
}

func Test_UpdateRequest_Email(t *testing.T) {
	t.Parallel()

	email := option.Some(RandomEmailAddress(t))
	req := &UpdateRequest{email: email}

	assert.Equal(t, email, req.Email())
}

func Test_UpdateRequest_PasswordHash(t *testing.T) {
	t.Parallel()

	passwordHash := option.Some(RandomPasswordHash(t))
	req := &UpdateRequest{passwordHash: passwordHash}

	assert.Equal(t, passwordHash, req.PasswordHash())
}

func Test_UpdateRequest_Username(t *testing.T) {
	t.Parallel()

	username := option.Some(RandomUsername(t))
	req := &UpdateRequest{username: username}

	assert.Equal(t, username, req.Username())
}

func Test_UpdateRequest_Role(t *testing.T) {
	t.Parallel()

	role := option.Some(RandomRole(t))
	req := &UpdateRequest{role: role}

	assert.Equal(t, role, req.Role())
}

func Test_UpdateRequest_GoString(t *testing.T) {
	t.Parallel()

	req := RandomUpdateRequest(t)

	expected := fmt.Sprintf(
		"UpdateRequest{userID:%#v, eTag:%#v, email:%#v, passwordHash:%#v, username:%#v, role:%#v,}",
		req.userID,
		req.eTag,
		req.email,
		req.passwordHash,
		req.username,
		req.role,
	)

	assert.Equal(t, expected, req.GoString())
}

func Test_UpdateRequest_String(t *testing.T) {
	t.Parallel()

	req := RandomUpdateRequest(t)

	expected := fmt.Sprintf(
		"{%s %s %s %s %s %s}",
		req.userID,
		req.eTag,
		req.email,
		req.passwordHash,
		req.username,
		req.role,
	)

	assert.Equal(t, expected, req.String())
}

func Test_UpdateRequest_Equal(t *testing.T) {
	t.Parallel()

	password := RandomPasswordCandidate()
	passwordHash, err := ParsePassword(password)
	require.NoError(t, err)

	userID := uuid.New()
	eTag := etag.Random()
	email := option.Some(RandomEmailAddress(t))
	username := option.Some(RandomUsername(t))
	role := option.Some(RandomRole(t))

	req1 := NewUpdateRequest(
		userID,
		eTag,
		email,
		option.Some(passwordHash),
		username,
		role,
	)
	req2 := NewUpdateRequest(
		userID,
		eTag,
		email,
		option.Some(passwordHash),
		username,
		role,
	)
	req3 := RandomUpdateRequest(t)

	testCases := []struct {
		name     string
		req1     *UpdateRequest
		req2     *UpdateRequest
		password option.Option[string]
		want     bool
	}{
		{
			name:     "equal requests",
			req1:     req1,
			req2:     req2,
			password: option.Some(password),
			want:     true,
		},
		{
			name:     "different requests",
			req1:     req1,
			req2:     req3,
			password: option.Some(password),
			want:     false,
		},
		{
			name:     "wrong password",
			req1:     req1,
			req2:     req2,
			password: option.Some("wrong_password"),
			want:     false,
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			got := tc.req1.Equal(tc.req2, tc.password)
			assert.Equal(t, tc.want, got)
		})
	}
}
