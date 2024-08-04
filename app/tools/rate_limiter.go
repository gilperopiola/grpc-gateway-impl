package tools

import (
	"github.com/gilperopiola/grpc-gateway-impl/app/core"

	"golang.org/x/time/rate"
)

var _ core.RateLimiter = &rateLimiter{}

type rateLimiter struct {
	*rate.Limiter
}

func NewRateLimiter(cfg *core.RLimiterCfg) core.RateLimiter {
	limit := rate.Limit(cfg.TokensPerSecond)
	limiter := rate.NewLimiter(limit, cfg.MaxTokens)
	return &rateLimiter{limiter}
}

/* -~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~- */

// True if it's allowed.
// False if it's not allowed.
//
// I know this seems rather useless, but it's a building block that can be further developed.
func (rl *rateLimiter) AllowRate() bool {
	return rl.Allow()
}
