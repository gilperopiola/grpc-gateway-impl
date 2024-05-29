package toolbox

import (
	"errors"

	"github.com/gilperopiola/god"
	"github.com/gilperopiola/grpc-gateway-impl/app/core"
	"github.com/gilperopiola/grpc-gateway-impl/app/core/errs"

	"golang.org/x/time/rate"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var _ core.RateLimiter = (*rateLimiter)(nil)

type rateLimiter struct {
	*rate.Limiter
}

func NewRateLimiter(cfg *core.RLimiterCfg) core.RateLimiter {
	limit := rate.Limit(cfg.TokensPerSecond)
	limiter := rate.NewLimiter(limit, cfg.MaxTokens)
	return &rateLimiter{limiter}
}

/* -~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~- */

func (rl *rateLimiter) LimitGRPC(ctx god.Ctx, req any, _ *god.GRPCInfo, handler god.GRPCHandler) (any, error) {
	if !rl.Allow() {
		core.LogUnexpectedErr(errors.New("rate limit exceeded"))
		return nil, status.Errorf(codes.ResourceExhausted, errs.RateLimitedMsg)
	}
	return handler(ctx, req)
}

/* -~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~- */
