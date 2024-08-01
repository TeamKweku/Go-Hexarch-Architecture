package user

import (
	"math/rand"
	"testing"
	"time"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"github.com/teamkweku/code-odessey-hex-arch/pkg/etag"
	"github.com/teamkweku/code-odessey-hex-arch/pkg/option"
)

func RandomEmailAddressCandidate() string {
	return gofakeit.Email()
}

func RandomUsernameCandidate() string {
	return gofakeit.Regex(usernamePattern)
}

func RandomPasswordCandidate() string {
	length := rand.Intn(PasswordMaxLen-PasswordMinLen) + PasswordMinLen
	raw := gofakeit.Password(true, true, true, true, true, length)
	return raw
}

func RandomRoleCandidate() int {
	roles := []Role{RoleReader, RoleAuthor, RoleEditor, RoleAdmin}
	return int(roles[rand.Intn(len(roles))])
}

func RandomEmailAddress(t *testing.T) EmailAddress {
	t.Helper()

	email, err := ParseEmailAddress(RandomEmailAddressCandidate())
	require.NoError(t, err)

	return email
}

func RandomUsername(t *testing.T) Username {
	t.Helper()

	username, err := ParseUsername(RandomUsernameCandidate())
	require.NoError(t, err)

	return username
}

func RandomPasswordHash(t *testing.T) PasswordHash {
	t.Helper()

	password := RandomPasswordCandidate()
	hash, err := ParsePassword(password)
	require.NoError(t, err)

	return hash
}

func RandomRole(t *testing.T) Role {
	t.Helper()

	role, err := ParseRole(RandomRoleCandidate())
	require.NoError(t, err)

	return role
}

func RandomOption[T any](t *testing.T) option.Option[T] {
	t.Helper()

	if rand.Intn(2) == 0 {
		switch any(*new(T)).(type) {
		case EmailAddress:
			email := any(RandomEmailAddress(t)).(T)
			return option.Some(email)
		case Username:
			username := any(RandomUsername(t)).(T)
			return option.Some(username)
		case PasswordHash:
			password := any(RandomPasswordHash(t)).(T)
			return option.Some(password)
		case Role:
			role := any(RandomRole(t)).(T)
			return option.Some(role)
		default:
			require.FailNow(
				t,
				"Unsupported type passed to RandomOption",
				"RandomOption does not support type %T",
				any(*new(T)),
			)
		}
	}

	return option.None[T]()
}

func RandomOptionFromInstance[T any](instance T) option.Option[T] {
	if rand.Intn(2) == 0 {
		return option.Some(instance)
	}

	return option.None[T]()
}

func RandomRegistrationRequest(t *testing.T) *RegistrationRequest {
	t.Helper()

	username := RandomUsername(t)
	email := RandomEmailAddress(t)
	password := RandomPasswordHash(t)
	return NewRegistrationRequest(username, email, password)

}

func RandomLoginRequest(t *testing.T) *LoginRequest {
	t.Helper()

	email := RandomEmailAddress(t)
	passwordCandidate := RandomPasswordCandidate()
	return NewLoginRequest(email, passwordCandidate)
}

func RandomUpdateRequest(t *testing.T) *UpdateRequest {
	t.Helper()

	id := uuid.New()
	eTag := etag.Random()
	username := RandomOption[Username](t)
	email := RandomOption[EmailAddress](t)
	password := RandomOption[PasswordHash](t)
	role := RandomOption[Role](t)

	return NewUpdateRequest(id, eTag, email, password, username, role)

}

func RandomUser(t *testing.T) *User {
	t.Helper()

	id := uuid.New()
	etag := etag.Random()
	username := RandomUsername(t)
	email := RandomEmailAddress(t)
	password := RandomPasswordHash(t)
	role := RandomRole(t)
	createdAt := time.Now()
	passwordUpdatedAt := time.Now()

	return NewUser(id, etag, username, email, password, role, createdAt, passwordUpdatedAt)
}
