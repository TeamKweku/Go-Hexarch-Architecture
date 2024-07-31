package user

import (
	"fmt"
	"net/mail"
	"regexp"

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
	PasswordMinLen = 0
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
type passwordComparator func(hash PasswordHash, candidate string) error

func bcryptCompare(hash PasswordHash, candidate string) error {
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
// Direct comparison of password hashes is impossible by design.
func (r *RegistrationRequest) Equal(other *RegistrationRequest, password string) bool {
	if len(r.passwordHash.bytes) > 0 || len(other.passwordHash.bytes) > 0 {
		if err := bcryptCompare(r.passwordHash, password); err != nil {
			return false
		}
		if err := bcryptCompare(other.passwordHash, password); err != nil {
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
