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
				assert.NoError(t, bcryptCompare(hash, candidate))
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
		t.Run(tt.name, func(t *testing.T) {
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
		t.Run(tt.expected, func(t *testing.T) {
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
		t.Run(fmt.Sprintf("Role(%d)", int(tt.role)), func(t *testing.T) {
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

	user := NewUser(id, eTag, username, email, passwordHash, role, createdAt, passwordChangedAt)

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

	expected := fmt.Sprintf("{ %s %s %s %s %s %s}",
		user.id, user.eTag, user.username, user.email, user.passwordHash, user.role)

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

			req, err := ParseRegistrationRequest(tc.username, tc.email, tc.password)

			if tc.wantError {
				assert.Error(t, err)
				assert.Nil(t, req)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, req)
				assert.Equal(t, tc.username, req.Username().String())
				assert.Equal(t, tc.email, req.Email().String())
				assert.NoError(t, bcryptCompare(req.PasswordHash(), tc.password))
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
	req3 := NewRegistrationRequest(req1.Username(), req1.Email(), req1.PasswordHash())

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
