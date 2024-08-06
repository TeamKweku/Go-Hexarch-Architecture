package user

import (
	"fmt"
	"net/mail"
	"regexp"
	"time"

	"github.com/google/uuid"
	"github.com/teamkweku/code-odessey-hex-arch/pkg/etag"
	"github.com/teamkweku/code-odessey-hex-arch/pkg/option"
	"golang.org/x/crypto/bcrypt"
)

// EmailAddress is a dedicated usernameCandidate type for valid email addresses. New
// instances are validated for RFC5332 compliance.
type EmailAddress struct {
	raw string
}

// ParseEmailAddress returns a new email address from `candidate`, validating
// that the email address conforms to RFC5332 standards (with the minor
// divergences introduce by the Go standard library, documented in [net/mail]).
func ParseEmailAddress(candidate string) (EmailAddress, error) {
	if _, err := mail.ParseAddress(candidate); err != nil {
		return EmailAddress{}, NewEmailAddressFormatError(candidate)
	}

	return EmailAddress{raw: candidate}, nil
}

// String returns the raw email address.
func (ea EmailAddress) String() string {
	return ea.raw
}

const (
	// UsernameMinLen is the minimum length of a username in bytes.
	UsernameMinLen = 3

	// UsernameMaxLen is the maximum length of a username in bytes.
	UsernameMaxLen = 16

	UsernamePatternTemplate = "^[a-zA-Z0-9_]{%d,%d}$"
)

var (
	// usernamePattern -> "^[a-zA-Z0-9_]{3,16}$"
	usernamePattern = fmt.Sprintf(UsernamePatternTemplate, UsernameMinLen, UsernameMaxLen)
	usernameRegex   = regexp.MustCompile(usernamePattern)
)

// Username represents a valid Username.
type Username struct {
	raw string
}

// ParseUsername returns either a valid [Username] or an error indicating why
// the raw username was invalid.
func ParseUsername(candidate string) (Username, error) {
	if len(candidate) < UsernameMinLen {
		return Username{}, NewUsernameTooShortError()
	}
	if len(candidate) > UsernameMaxLen {
		return Username{}, NewUsernameTooLongError()
	}
	if !usernameRegex.MatchString(candidate) {
		return Username{}, NewUsernameFormatError()
	}

	return Username{raw: candidate}, nil
}

func (u Username) String() string {
	return u.raw
}

const (
	PasswordMinLen = 8
	PasswordMaxLen = 72
)

// PasswordHash represents a validated and hashed password.
//
// The hash is obfuscated when printed with the %s, %v and %#v verbs.
//
// CAUTION: The fmt package uses reflection to print unexported fields without
// invoking their String or GoString methods. Printing structs containing
// unexported PasswordHashes will result in the hash bytes being exposed.
type PasswordHash struct {
	bytes []byte
}

func ParsePassword(candidate string) (PasswordHash, error) {
	return parsePassword(candidate, bcryptHash)
}

func parsePassword(candidate string, hasher passwordHasher) (PasswordHash, error) {
	if err := validatePasswordCandidate(candidate); err != nil {
		return PasswordHash{}, err
	}

	hash, err := hasher(candidate)
	if err != nil {
		return PasswordHash{}, err
	}

	return hash, nil
}

func validatePasswordCandidate(candidate string) error {
	if len(candidate) < PasswordMinLen {
		return NewPasswordTooShortError()
	}
	if len(candidate) > PasswordMaxLen {
		return NewPasswordTooLongError()
	}
	return nil
}

// NewPasswordHashFromTrustedSource wraps a hashed password in a [PasswordHash].
func NewPasswordHashFromTrustedSource(raw []byte) PasswordHash {
	return PasswordHash{bytes: raw}
}

func (ph PasswordHash) Bytes() []byte {
	return ph.bytes
}

// String obfuscates the hash bytes when the hash is printed with the %s and %v
// verbs.
func (ph PasswordHash) String() string {
	return "{REDACTED}"
}

// GoString obfuscates the hash bytes when the hash is printed with the %#v verb.
func (ph PasswordHash) GoString() string {
	return "PasswordHash{bytes:REDACTED}"
}

// passwordHasher is a function that hashes a password candidate. By abstracting
// a general class of hasher functions, we can simulate hashing errors in tests.
type passwordHasher func(candidate string) (PasswordHash, error)

func bcryptHash(candidate string) (PasswordHash, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(candidate), bcrypt.DefaultCost)
	if err != nil {
		return PasswordHash{}, fmt.Errorf("hash password: %w", err)
	}
	return NewPasswordHashFromTrustedSource(hash), nil
}

// passwordComparator is a function that compares a [PasswordHash] and password.
// By abstracting a general class of comparator functions, we can simulate
// comparison errors in tests.
type PasswordComparator func(hash PasswordHash, candidate string) error

func BcryptCompare(hash PasswordHash, candidate string) error {
	if err := bcrypt.CompareHashAndPassword(hash.bytes, []byte(candidate)); err != nil {
		return &AuthError{Cause: err}
	}
	return nil
}

// Role specifies various roles for the blog application
type Role int

const (
	RoleReader Role = iota
	RoleAuthor
	RoleEditor
	RoleAdmin
)

// ParseRole returns a valid Role or an error
func ParseRole(candidate int) (Role, error) {
	if candidate < int(RoleReader) || candidate > int(RoleAdmin) {
		return Role(0), fmt.Errorf("invalid role: %d", candidate)
	}

	return Role(candidate), nil
}

// String returns the string representation of the Role
func (r Role) String() string {
	switch r {
	case RoleReader:
		return "Reader"
	case RoleAuthor:
		return "Author"
	case RoleEditor:
		return "Editor"
	case RoleAdmin:
		return "Admin"
	default:
		return "Unknown"
	}
}

// GoString returns a Go-syntax representation of the Role
func (r Role) GoString() string {
	return fmt.Sprintf("Role(%d)", int(r))
}

func NewUser(
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
		"User{id:%#v, eTag:%#v, username:%#v, email:%#v, passwordHash:%#v,",
		u.id, u.eTag, u.username, u.email, u.passwordHash,
	)
}

// GoString ensures that the [PasswordHash]'s GoString method is invoked when the
// User is printed with the %s or %v verbs. Unexported fields are otherwise printed
// reflectively, which would expose the hash.
func (u User) String() string {
	return fmt.Sprintf("{ %s %s %s %s %s %s}",
		u.id, u.eTag, u.username, u.email, u.passwordHash, u.role)
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

// ParseRegistrationRequest returns a new [RegistrationRequest] from raw inputs.
//
// # Errors
//   - [ValidationErrors], if one or more inputs are invalid.
//   - Unexpected internal response.
func ParseRegistrationRequest(
	usernameCandidate string,
	emailCandidate string,
	passwordCandidate string,
) (*RegistrationRequest, error) {
	var validationErrs ValidationErrors
	username, err := ParseUsername(usernameCandidate)
	if pushErr := validationErrs.PushValidationError(err); pushErr != nil {
		return nil, pushErr
	}

	email, err := ParseEmailAddress(emailCandidate)
	if pushErr := validationErrs.PushValidationError(err); pushErr != nil {
		return nil, pushErr
	}

	passwordHash, err := ParsePassword(passwordCandidate)
	if pushErr := validationErrs.PushValidationError(err); pushErr != nil {
		return nil, pushErr
	}

	if validationErrs.Any() {
		return nil, validationErrs
	}

	return NewRegistrationRequest(username, email, passwordHash), nil
}

// getter for RegistrationRequest
func (r *RegistrationRequest) Username() Username {
	return r.username
}

func (r *RegistrationRequest) Email() EmailAddress {
	return r.email
}

func (r *RegistrationRequest) PasswordHash() PasswordHash {
	return r.passwordHash
}

// Equal returns true if `r.passwordHash` can be obtained from `password`,
// and the two requests are equal in all other fields.
//
// Direct comparu.pkid,ison of password hashes is impossible by design.
func (r *RegistrationRequest) Equal(other *RegistrationRequest, password string) bool {
	if len(r.passwordHash.bytes) > 0 || len(other.passwordHash.bytes) > 0 {
		if err := BcryptCompare(r.passwordHash, password); err != nil {
			return false
		}
		if err := BcryptCompare(other.passwordHash, password); err != nil {
			return false
		}
	}

	return r.username == other.username && r.email == other.email
}

// GoString ensures that the [PasswordHash]'s GoString method is invoked when the
// request is printed with the %#v verb. Unexported fields are otherwise printed
// reflectively, which would expose the hash.
func (r RegistrationRequest) GoString() string {
	return fmt.Sprintf("RegistrationRequest{username:%#v, email:%#v, passwordHash:%#v}",
		r.username, r.email, r.passwordHash)
}

// String ensures that the [PasswordHash]'s String method is invoked when the
// request is printed with the %s or %v verbs. Unexported fields are otherwise
// printed reflectively, which would expose the hash.
func (r RegistrationRequest) String() string {
	return fmt.Sprintf("{%s %s %s}", r.username, r.email, r.passwordHash)
}

// AuthRequest describes the data required to authenticate a user.
type LoginRequest struct {
	email             EmailAddress
	passwordCandidate string
}

func NewLoginRequest(email EmailAddress, passwordCandidate string) *LoginRequest {
	return &LoginRequest{
		email:             email,
		passwordCandidate: passwordCandidate,
	}
}

// ParseLoginRequest returns a new [LoginRequest] from raw inputs.
//
// # Errors
// - [ValidationErrors], if `emailCandidate` is invalid
func ParseLoginRequest(emailCandidate string, passwordCandidate string) (*LoginRequest, error) {
	var validationErrs ValidationErrors
	email, err := ParseEmailAddress(emailCandidate)
	if pushErr := validationErrs.PushValidationError(err); pushErr != nil {
		return nil, pushErr
	}

	if validationErrs.Any() {
		return nil, validationErrs
	}

	return NewLoginRequest(email, passwordCandidate), nil
}

func (lg *LoginRequest) Email() EmailAddress {
	return lg.email
}

func (lg *LoginRequest) PasswordCandidate() string {
	return lg.passwordCandidate
}

// GoString ensures that `passwordCandidate` is obfuscated when the request is
// printed with the %#v verb. Unexported fields are otherwise printed
// reflectively, which would expose the hash.
func (lg LoginRequest) GoString() string {
	return fmt.Sprintf("LoginRequest{email:%#v, passwordCandidate:REDACTED}", lg.email)
}

// GoString ensures that `passwordCandidate` is obfuscated when the request is
// printed with the %s or %v verbs. Unexported fields are otherwise printed
// reflectively, which would expose the hash.
func (lg LoginRequest) String() string {
	return fmt.Sprintf("{%s REDACTED}", lg.email)
}

// UpdateRequest describes  the data fields required to update a user. All fields but
// userID are optional.
type UpdateRequest struct {
	userID       uuid.UUID
	eTag         etag.ETag
	username     option.Option[Username]
	email        option.Option[EmailAddress]
	passwordHash option.Option[PasswordHash]
	role         option.Option[Role]
}

func NewUpdateRequest(
	userID uuid.UUID,
	eTag etag.ETag,
	email option.Option[EmailAddress],
	passwordHash option.Option[PasswordHash],
	username option.Option[Username],
	role option.Option[Role],
) *UpdateRequest {
	return &UpdateRequest{
		userID:       userID,
		eTag:         eTag,
		email:        email,
		passwordHash: passwordHash,
		username:     username,
		role:         role,
	}
}

// ParseUpdateRequest returns a new [UpdateRequest] from raw inputs.
//
// # Errors
//   - [ValidationErrors], if one or more inputs are invalid.
//   - Unexpected internal response.
func ParseUpdateRequest(
	userID uuid.UUID,
	eTag etag.ETag,
	emailCandidate option.Option[string],
	passwordCandidate option.Option[string],
	usernameCandidate option.Option[string],
	roleCandidate option.Option[int],

) (*UpdateRequest, error) {
	var validationErrs ValidationErrors

	email, err := option.Map(emailCandidate, ParseEmailAddress)
	if pushErr := validationErrs.PushValidationError(err); pushErr != nil {
		return nil, pushErr
	}

	passwordHash, err := option.Map(passwordCandidate, ParsePassword)
	if pushErr := validationErrs.PushValidationError(err); pushErr != nil {
		return nil, pushErr
	}

	username, err := option.Map(usernameCandidate, ParseUsername)
	if pushErr := validationErrs.PushValidationError(err); pushErr != nil {
		return nil, pushErr
	}

	role, err := option.Map(roleCandidate, ParseRole)
	if pushErr := validationErrs.PushValidationError(err); pushErr != nil {
		return nil, pushErr
	}

	if validationErrs.Any() {
		return nil, validationErrs
	}

	return NewUpdateRequest(userID, eTag, email, passwordHash, username, role), nil
}

func (ur *UpdateRequest) UserID() uuid.UUID {
	return ur.userID
}

func (ur *UpdateRequest) ETag() etag.ETag {
	return ur.eTag
}

func (ur *UpdateRequest) Email() option.Option[EmailAddress] {
	return ur.email
}

func (ur *UpdateRequest) PasswordHash() option.Option[PasswordHash] {
	return ur.passwordHash
}

func (ur *UpdateRequest) Username() option.Option[Username] {
	return ur.username
}

func (ur *UpdateRequest) Role() option.Option[Role] {
	return ur.role
}

func (ur UpdateRequest) GoString() string {
	return fmt.Sprintf("UpdateRequest{userID:%#v, eTag:%#v, email:%#v, passwordHash:%#v, username:%#v, role:%#v,}",
		ur.userID, ur.eTag, ur.email, ur.passwordHash, ur.username, ur.role)
}

func (ur UpdateRequest) String() string {
	return fmt.Sprintf("{%s %s %s %s %s %s}", ur.userID, ur.eTag, ur.email, ur.passwordHash, ur.username, ur.role)
}

// Equal returns true if `ur.passwordHash` can be obtained from `password`,
// and the two requests are equal in all other fields.
//
// Direct comparison of password hashes is impossible by design.
func (ur *UpdateRequest) Equal(other *UpdateRequest, password option.Option[string]) bool {
	if ur.passwordHash.IsSome() || other.passwordHash.IsSome() {
		pw := password.UnwrapOrZero()
		if err := BcryptCompare(ur.passwordHash.UnwrapOrZero(), pw); err != nil {
			return false
		}
		if err := BcryptCompare(other.passwordHash.UnwrapOrZero(), pw); err != nil {
			return false
		}
	}

	if ur.userID != other.userID ||
		ur.email != other.email {
		return false
	}

	return true
}
