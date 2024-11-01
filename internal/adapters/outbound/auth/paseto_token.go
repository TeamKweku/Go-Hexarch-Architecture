package auth

import (
	"fmt"
	"time"

	"aidanwoods.dev/go-paseto"

	"github.com/teamkweku/code-odessey-hex-arch/internal/core/domain/auth"
	"github.com/teamkweku/code-odessey-hex-arch/internal/core/domain/user"
	"github.com/teamkweku/code-odessey-hex-arch/internal/core/ports/inbound"
)

/**
* PasetoToken implements the port inbound.TokenService interface
 */
type PasetoToken struct {
	token  *paseto.Token
	key    *paseto.V4SymmetricKey
	parser *paseto.Parser
}

// New creates a new paseto instance and returns an token interface
// just to make sure the function implements all methods of interface
// CreateToken & VerifyToken methods
func NewPasetoToken() (inbound.TokenService, error) {
	token := paseto.NewToken()
	key := paseto.NewV4SymmetricKey()
	parser := paseto.NewParser()

	return &PasetoToken{
		token:  &token,
		key:    &key,
		parser: &parser,
	}, nil
}

// CreateToken creates a new paseto token
func (pt *PasetoToken) CreateToken(
	user *user.User,
	durationStr string,
) (string, *auth.Payload, error) {
	payload, err := auth.NewPayload(user, durationStr)
	if err != nil {
		return "", nil, fmt.Errorf("failed to set token payload: %w", err)
	}

	err = pt.token.Set("payload", payload)
	if err != nil {
		return "", nil, fmt.Errorf(
			"failed to set token payload: %w",
			auth.NewTokenCreationError("token payload error"),
		)
	}

	pt.token.SetIssuedAt(payload.IssuedAt)
	pt.token.SetNotBefore(payload.IssuedAt)
	pt.token.SetExpiration(payload.ExpiredAt)

	token := pt.token.V4Encrypt(*pt.key, nil)

	return token, payload, nil
}

func (pt *PasetoToken) VerifyToken(tokenstring string) (*auth.Payload, error) {
	var payload *auth.Payload

	parsedToken, err := pt.parser.ParseV4Local(*pt.key, tokenstring, nil)
	if err != nil {
		if err.Error() == "this token has expired" {
			return nil, fmt.Errorf("token has expired: %w", auth.NewExpiredTokenError(time.Now()))
		}
		return nil, fmt.Errorf("invalid token: %w", auth.NewInvalidTokenError("token error"))
	}

	err = parsedToken.Get("payload", &payload)
	if err != nil {
		return nil, fmt.Errorf(
			"invalid token payload %w",
			auth.NewInvalidTokenError("payload error"),
		)
	}

	// Normalize time zones to UTC
	payload.IssuedAt = payload.IssuedAt.UTC()
	payload.ExpiredAt = payload.ExpiredAt.UTC()

	return payload, nil
}
