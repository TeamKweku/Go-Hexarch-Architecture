// This file contains code inspired by https://github.com/angusgmorrison/realworld-go
// Original author: Angus Morrison
package user

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/teamkweku/code-odessey-hex-arch/pkg/etag"
)

// User is the central domain type for this package.
type User struct {
	pkid              int64
	id                uuid.UUID
	eTag              etag.ETag
	username          Username
	email             EmailAddress
	passwordHash      PasswordHash
	role              Role
	createAt          time.Time
	passwordChangedAt time.Time
}

func NewUser(
	pikd int64,
	id uuid.UUID,
	eTag etag.ETag,
	username Username,
	email EmailAddress,
	passwordHash PasswordHash,
	role Role,
	createAt time.Time,
	paswordChangedAt time.Time,

) *User {
	return &User{
		pkid:              pikd,
		id:                id,
		eTag:              eTag,
		username:          username,
		email:             email,
		passwordHash:      passwordHash,
		role:              role,
		createAt:          createAt,
		passwordChangedAt: paswordChangedAt,
	}
}

// getter methods
func (u *User) PKID() int64 {
	return u.pkid
}

func (u *User) ID() uuid.UUID {
	return u.id
}

func (u *User) ETag() etag.ETag {
	return u.eTag
}

func (u *User) Username() Username {
	return u.username
}

func (u *User) Email() EmailAddress {
	return u.email
}

func (u *User) PasswordHash() PasswordHash {
	return u.passwordHash
}

// GoString ensures that the [PasswordHash]'s GoString method is invoked when the
// User is printed with the %#v verb. Unexported fields are otherwise printed
// reflectively, which would expose the hash.
func (u User) GoString() string {
	return fmt.Sprintf(
		"User{pkid:%#v ,id:%#v, eTag:%#v, username:%#v, email:%#v, passwordHash:%#v,",
		u.pkid, u.id, u.eTag, u.username, u.email, u.passwordHash,
	)
}

// GoString ensures that the [PasswordHash]'s GoString method is invoked when the
// User is printed with the %s or %v verbs. Unexported fields are otherwise printed
// reflectively, which would expose the hash.
func (u User) String() string {
	return fmt.Sprintf("{%d, %s %s %s %s %s %s}",
		u.pkid, u.id, u.eTag, u.username, u.email, u.passwordHash, u.role)
}

// RegistrationRequest carries validated data required to register a new user.
type RegistrationRequest struct {
	username     Username
	email        EmailAddress
	passwordHash PasswordHash
}

func NewRegistrationRequest(
	username Username, email EmailAddress, passwordHash PasswordHash,
) *RegistrationRequest {
	return &RegistrationRequest{
		username:     username,
		email:        email,
		passwordHash: passwordHash,
	}
}
