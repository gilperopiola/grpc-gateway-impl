package service

import (
	"context"
	"fmt"

	"github.com/gilperopiola/grpc-gateway-impl/app/core"
	"github.com/gilperopiola/grpc-gateway-impl/app/core/pbs"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type HealthSubService struct {
	pbs.UnimplementedHealthServiceServer
	Tools core.Tools
}

// Checks the DB connection and makes an HTTP call.
// If both succeed, returns OK and sends the app version.
func (h *HealthSubService) CheckHealth(ctx context.Context, _ *pbs.CheckHealthRequest) (*pbs.CheckHealthResponse, error) {

	// Get the DB or return unhealthy.
	if h.Tools.GetDB() == nil {
		return nil, status.Error(codes.Unavailable, "database connection unhealthy")
	}

	// Make HTTP call or return unhealthy.
	if _, err := h.Tools.GetCurrentWeather(ctx, 50, 50); err != nil {
		core.LogUnexpected(fmt.Errorf("network call unhealthy: %w", err))
		return nil, status.Error(codes.Unavailable, "network call unhealthy")
	}

	msg := core.AppName + " " + core.AppVersion + " healthy"
	return &pbs.CheckHealthResponse{Info: msg}, nil
}
