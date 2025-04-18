package service

import (
	"context"
	"fmt"
	"strings"

	"github.com/gilperopiola/grpc-gateway-impl/app/core"
	"github.com/gilperopiola/grpc-gateway-impl/app/core/pbs"
	"github.com/gilperopiola/grpc-gateway-impl/app/core/utils"
	"go.uber.org/zap"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type HealthSvc struct {
	pbs.UnimplementedHealthServiceServer
	Clients core.Clients
	Tools   core.Tools
}

// Checks the DB connection and makes an HTTP call.
// If both succeed, returns OK and sends the app version.
func (h *HealthSvc) CheckHealth(ctx context.Context, _ *pbs.CheckHealthRequest) (*pbs.CheckHealthResponse, error) {

	msg := core.G.AppName + " " + core.G.Version

	// Get the DB or return unhealthy.
	if _, err := utils.TryAndRetryNoErr(h.Clients.GetDB, 3, false, nil); err != nil {
		return nil, status.Error(codes.Unavailable, msg+" unhealthy: database connection not working")
	}

	// Make HTTP call or return unhealthy.
	word := "healthy"
	gptResponse, err := h.Clients.SendRequestToGPT(ctx, fmt.Sprintf("Give a really short response, which includes the word '%s'.", word))
	if err != nil {
		return nil, status.Error(codes.Unavailable, msg+" unhealthy: http calls not working")
	}

	if !strings.Contains(strings.ToLower(gptResponse), "healthy") {
		zap.S().Warnf("GPT response does not contain 'healthy': %s", gptResponse)
	}

	return &pbs.CheckHealthResponse{Info: msg + " healthy"}, nil
}
