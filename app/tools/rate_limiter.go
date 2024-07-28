package tools

import (
	"errors"

	"github.com/gilperopiola/god"
	"github.com/gilperopiola/grpc-gateway-impl/app/core"
	"github.com/gilperopiola/grpc-gateway-impl/app/core/errs"

	"golang.org/x/time/rate"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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

func (rl *rateLimiter) LimitGRPC(ctx god.Ctx, req any, _ *god.GRPCInfo, handler god.GRPCHandler) (any, error) {
	if !rl.Allow() {
		core.LogUnexpected(errors.New("rate limit exceeded"))
		return nil, status.Errorf(codes.ResourceExhausted, errs.RateLimitedMsg)
	}
	return handler(ctx, req)
}
