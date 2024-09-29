package service

import (
	"context"
	"strings"

	"github.com/gilperopiola/grpc-gateway-impl/app/core"
	"github.com/gilperopiola/grpc-gateway-impl/app/core/pbs"
	"github.com/gilperopiola/grpc-gateway-impl/app/core/utils"
	"go.uber.org/zap"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type HealthSubService struct {
	pbs.UnimplementedHealthServiceServer
	Clients core.Clients
	Tools   core.Tools
}

// Checks the DB connection and makes an HTTP call.
// If both succeed, returns OK and sends the app version.
func (h *HealthSubService) CheckHealth(ctx context.Context, _ *pbs.CheckHealthRequest) (*pbs.CheckHealthResponse, error) {

	msg := core.G.AppName + " " + core.G.Version

	// Get the DB or return unhealthy.
	if _, err := utils.RetryV2(h.Clients.GetDB, utils.BasicRetryCfg(2, nil)); err != nil {
		return nil, status.Error(codes.Unavailable, msg+" unhealthy: database connection not working")
	}

	// Make HTTP call or return unhealthy.
	if _, err := h.Clients.GetCurrentWeather(ctx, 50, 50); err != nil {
		gptResponse, err := h.Clients.SendToGPT(ctx, "give a really short response, which includes the word 'healthy'.")
		if err != nil {
			return nil, status.Error(codes.Unavailable, msg+" unhealthy: http calls not working")
		}

		if !strings.Contains(strings.ToLower(gptResponse), "healthy") {
			zap.S().Warnf("GPT response does not contain 'healthy': %s", gptResponse)
		}
	}

	msg += " healthy"

	return &pbs.CheckHealthResponse{Info: msg}, nil
}
