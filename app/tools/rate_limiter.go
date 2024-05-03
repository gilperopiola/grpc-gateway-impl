package tools

import (
	"context"
	"errors"

	"github.com/gilperopiola/grpc-gateway-impl/app/core"
	"github.com/gilperopiola/grpc-gateway-impl/app/core/errs"

	"golang.org/x/time/rate"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type rateLimiter struct {
	*rate.Limiter
}

func NewRateLimiter(cfg *core.RLimiterCfg) core.RateLimiter {
	limit := rate.Limit(cfg.TokensPerSecond)
	limiter := rate.NewLimiter(limit, cfg.MaxTokens)
	return &rateLimiter{limiter}
}

/* -~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~- */

func (rl *rateLimiter) LimitGRPC(ctx context.Context, req any, _ *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
	if !rl.Allow() {
		core.LogUnexpectedErr(errors.New("Rate limit exceeded!"))
		return nil, status.Errorf(codes.ResourceExhausted, errs.RateLimitedMsg)
	}
	return handler(ctx, req)
}
