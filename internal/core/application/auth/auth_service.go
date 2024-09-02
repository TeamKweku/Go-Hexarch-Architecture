package auth

import (
	"fmt"

	"github.com/teamkweku/code-odessey-hex-arch/config"
	"github.com/teamkweku/code-odessey-hex-arch/internal/core/domain/auth"
	"github.com/teamkweku/code-odessey-hex-arch/internal/core/domain/user"
	"github.com/teamkweku/code-odessey-hex-arch/internal/core/ports/inbound"
)

type AuthService struct {
	tokenService inbound.TokenService
}

func NewAuthService(tokenService inbound.TokenService) *AuthService {
	return &AuthService{tokenService}
}

func (as *AuthService) CreateToken(
	user *user.User,
	cfg config.Config,
) (string, *auth.Payload, error) {
	token, payload, err := as.tokenService.CreateToken(user, cfg)
	if err != nil {
		return "", nil, fmt.Errorf("auth service - create token: %w", err)
	}

	return token, payload, nil
}

func (as *AuthService) VerifyToken(token string) (*auth.Payload, error) {
	payload, err := as.tokenService.VerifyToken(token)
	if err != nil {
		return nil, fmt.Errorf("auth service - verify token: %w", err)
	}

	return payload, nil
}
