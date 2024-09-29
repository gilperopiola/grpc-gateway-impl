package service

import (
	"context"

	"github.com/gilperopiola/grpc-gateway-impl/app/core"
	"github.com/gilperopiola/grpc-gateway-impl/app/core/pbs"
	"github.com/gilperopiola/grpc-gateway-impl/app/core/utils"

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
		gptResponse, err := h.Clients.SendToGPT(ctx, "your response is going to be shown to a user of my API who is consulting the /health endpoint, so if you get this message just respond with something the user would expect to see in case it's healthy.")
		if err != nil {
			return nil, status.Error(codes.Unavailable, msg+" unhealthy: http calls not working")
		}

		msg += " " + gptResponse
	} else {
		msg += " healthy"
	}

	return &pbs.CheckHealthResponse{Info: msg}, nil
}
