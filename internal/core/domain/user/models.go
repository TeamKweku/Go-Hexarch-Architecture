package user

import (
	"time"

	"github.com/google/uuid"

	"github.com/teamkweku/code-odessey-hex-arch/pkg/etag"
)

// User is the central domain type for this package.
type User struct {
	id                uuid.UUID
	eTag              etag.ETag
	username          Username
	email             EmailAddress
	passwordHash      PasswordHash
	role              Role
	createAt          time.Time
	passwordChangedAt time.Time
}
