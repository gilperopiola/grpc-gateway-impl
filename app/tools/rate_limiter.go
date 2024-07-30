package tools

import (
	"github.com/gilperopiola/grpc-gateway-impl/app/core"

	"golang.org/x/time/rate"
)

type rateLimiter struct {
	*rate.Limiter
}

func NewRateLimiter(cfg *core.RLimiterCfg) *rateLimiter {
	limit := rate.Limit(cfg.TokensPerSecond)
	limiter := rate.NewLimiter(limit, cfg.MaxTokens)
	return &rateLimiter{limiter}
}

var _ core.RateLimiter = (*rateLimiter)(nil)

/* -~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~- */

// Returns true if it's allowed.
// Returns false if it's not allowed.
//
// I know this seems rather useless, but it's a building block that can be further developed.
func (rl *rateLimiter) AllowRate() bool {
	return rl.Allow()
}
