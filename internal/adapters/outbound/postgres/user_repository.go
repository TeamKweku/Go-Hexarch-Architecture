package postgres

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgtype"

	"github.com/teamkweku/code-odessey-hex-arch/internal/adapters/outbound/postgres/sqlc"
	"github.com/teamkweku/code-odessey-hex-arch/internal/core/domain/user"
	"github.com/teamkweku/code-odessey-hex-arch/internal/core/ports/outbound"
	"github.com/teamkweku/code-odessey-hex-arch/pkg/etag"
)

var _ outbound.UserRepository = (*Client)(nil)

// GetUserByID returns the [user.User] with the given ID or
// [user.NotFoundErr] if no such user is found
func (c *Client) GetUserByID(
	ctx context.Context,
	id uuid.UUID,
) (*user.User, error) {
	return getUserById(ctx, c.queries, id)
}

func getUserById(
	ctx context.Context,
	q queries,
	id uuid.UUID,
) (*user.User, error) {
	row, err := q.GetUserById(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, user.NewNotFoundByIDError(id)
		}
		return nil, fmt.Errorf("database error: %w", err)
	}

	return parseUser(
		row.ID,
		row.Username,
		row.Email,
		row.PasswordHash,
		row.Role,
		row.PasswordChangedAt,
		row.CreatedAt,
		row.UpdatedAt,
	)
}

// GetUserByEmail returns the [user.User] with the given email, or
// [user.NotFoundErr] if no user exists with email
func (c *Client) GetUserByEmail(
	ctx context.Context,
	email user.EmailAddress,
) (*user.User, error) {
	return getUserByEmail(ctx, c.queries, email)
}

func getUserByEmail(
	ctx context.Context,
	q queries,
	email user.EmailAddress,
) (*user.User, error) {
	row, err := q.GetUserByEmail(ctx, email.String())
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, user.NewNotFoundByEmailError(email)
		}
		return nil, fmt.Errorf("database error: %w", err)
	}

	return parseUser(
		row.ID,
		row.Username,
		row.Email,
		row.PasswordHash,
		row.Role,
		row.PasswordChangedAt,
		row.CreatedAt,
		row.UpdatedAt,
	)
}

// CreateUser creates a new user record from the given
// [user.RegistrationRequest] and returns the created [user.User].
//
// Returns [user.ValidationError] if database constraints are violated.
func (c *Client) CreateUser(
	ctx context.Context,
	req *user.RegistrationRequest,
) (*user.User, error) {
	return createUser(ctx, c.queries, req)
}

func createUser(
	ctx context.Context,
	q queries,
	req *user.RegistrationRequest,
) (*user.User, error) {
	row, err := q.CreateUser(ctx, newCreateUserParamsFromRegistrationReq(req))
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			return nil, createUserErrToDomain(pgErr, req)
		}
		return nil, fmt.Errorf(
			"create user record from request %#v: %w",
			req,
			err,
		)
	}

	return parseUser(
		row.ID,
		row.Username,
		row.Email,
		row.PasswordHash,
		row.Role,
		row.PasswordChangedAt,
		row.CreatedAt,
		row.UpdatedAt,
	)
}

func newCreateUserParamsFromRegistrationReq(
	req *user.RegistrationRequest,
) sqlc.CreateUserParams {
	return sqlc.CreateUserParams{
		Email:        req.Email().String(),
		Username:     req.Username().String(),
		PasswordHash: string(req.PasswordHash().Bytes()),
	}
}

func createUserErrToDomain(
	pgErr *pgconn.PgError,
	req *user.RegistrationRequest,
) error {
	switch pgErr.Code {
	case "23505": // unique_violation
		switch pgErr.ConstraintName {
		case "users_username_key":
			return user.NewDuplicateUsernameError(req.Username())
		case "users_email_key":
			return user.NewDuplicateEmailError(req.Email())
		}
	}
	return fmt.Errorf("unexpected database error: %w", pgErr)
}

// UpdateUser updates the user record and returns the updated [user.User]
// Returns [user.ValidationError] if database constriants are violated
func (c *Client) UpdateUser(
	ctx context.Context,
	req *user.UpdateRequest,
) (*user.User, error) {
	tx, err := c.db.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("begin UpdateUser transaction: %w", err)
	}

	defer func() { _ = tx.Rollback(ctx) }()

	queries := sqlc.New(tx)
	usr, err := updateUser(ctx, queries, req)
	if err != nil {
		return nil, err
	}

	if err := tx.Commit(ctx); err != nil {
		// TODO: Handle serialization failures
		return nil, fmt.Errorf("commit UpdateUser transaction: %w", err)
	}

	return usr, nil
}

func updateUser(
	ctx context.Context,
	q queries,
	req *user.UpdateRequest,
) (*user.User, error) {
	ok, err := q.UserExists(ctx, req.UserID())
	if err != nil {
		return nil, fmt.Errorf(
			"query existence of user with ID %q: %w",
			req.UserID().String(), err)
	}
	if !ok {
		return nil, user.NewNotFoundByIDError(req.UserID())
	}

	params := parseUpdateUserParams(req)
	row, err := q.UpdateUser(ctx, params)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, &user.ConcurrentModificationError{
				ID:   req.UserID(),
				ETag: req.ETag(),
			}
		}

		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			if strings.Contains(pgErr.ConstraintName, "email") {
				return nil, user.NewDuplicateEmailError(
					req.Email().UnwrapOrZero(),
				)
			}
		}

		return nil, fmt.Errorf("database error: %w", err)
	}

	return parseUser(
		row.ID,
		row.Username,
		row.Email,
		row.PasswordHash,
		row.Role,
		row.PasswordChangedAt,
		row.CreatedAt,
		row.UpdatedAt,
	)
}

func parseUpdateUserParams(req *user.UpdateRequest) sqlc.UpdateUserParams {
	params := sqlc.UpdateUserParams{
		ID: req.UserID(),
	}

	if req.Username().IsSome() {
		params.Username = pgtype.Text{
			String: req.Username().UnwrapOrZero().String(),
			Valid:  true,
		}
	}

	if req.Email().IsSome() {
		params.Email = pgtype.Text{
			String: req.Email().UnwrapOrZero().String(),
			Valid:  true,
		}
	}

	if req.PasswordHash().IsSome() {
		params.PasswordHash = pgtype.Text{
			String: string(req.PasswordHash().UnwrapOrZero().Bytes()),
			Valid:  true,
		}
	}

	if req.Role().IsSome() {
		params.Role = pgtype.Text{
			String: req.Role().UnwrapOrZero().String(),
			Valid:  true,
		}
	}

	return params
}

func parseUser(
	id uuid.UUID,
	username string,
	email string,
	passwordHash string,
	role string,
	passwordChangedAt time.Time,
	createdAt time.Time,
	updatedAt time.Time,
) (*user.User, error) {
	parsedID, err := uuid.Parse(id.String())
	if err != nil {
		return nil, fmt.Errorf(
			"userID %q from database can't be parsed: %w",
			id,
			err,
		)
	}

	eTag := etag.New(parsedID, updatedAt)

	parsedEmail, err := user.ParseEmailAddress(email)
	if err != nil {
		return nil, err
	}

	parsedUsername, err := user.ParseUsername(username)
	if err != nil {
		return nil, err
	}

	parsedPasswordHash := user.NewPasswordHashFromTrustedSource(
		[]byte(passwordHash),
	)

	roleInt, err := roleStringToInt(role)
	if err != nil {
		return nil, fmt.Errorf("invalid role from database: %w", err)
	}

	parsedRole, err := user.ParseRole(roleInt)
	if err != nil {
		return nil, err
	}

	return user.NewUser(
		parsedID,
		eTag,
		parsedUsername,
		parsedEmail,
		parsedPasswordHash,
		parsedRole,
		createdAt,
		passwordChangedAt,
		updatedAt,
	), nil
}

func roleStringToInt(roleStr string) (int, error) {
	switch strings.ToLower(roleStr) {
	case "reader":
		return int(user.RoleReader), nil
	case "admin":
		return int(user.RoleAdmin), nil
	default:
		return -1, fmt.Errorf("unknown role: %s", roleStr)
	}
}
