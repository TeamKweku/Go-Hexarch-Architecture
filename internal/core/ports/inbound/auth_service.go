package inbound

import (
	"github.com/teamkweku/code-odessey-hex-arch/internal/core/domain/auth"
	"github.com/teamkweku/code-odessey-hex-arch/internal/core/domain/user"
)

type TokenService interface {
	CreateToken(user *user.User, durationStr string) (string, *auth.Payload, error)
	VerifyToken(token string) (*auth.Payload, error)
}
