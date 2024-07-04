package toolbox

import (
	"context"
	"errors"

	"github.com/gilperopiola/grpc-gateway-impl/app/core"
	"github.com/gilperopiola/grpc-gateway-impl/app/core/pbs"
)

var _ core.HealthChecker = (*healthChecker)(nil)

type healthChecker struct {
	healthCheckFn ServiceFn
}

func NewHealthChecker(serviceFn ServiceFn) *healthChecker {
	return &healthChecker{serviceFn}
}

func (hc *healthChecker) CheckHealth() error {
	resp, err := hc.healthCheckFn(context.Background(), nil)
	if err != nil {
		return err
	}
	if resp == nil {
		return errors.New("health check failed")
	}
	return nil
}

type ServiceFn func(context.Context, *pbs.AnswerGroupInviteRequest) (*pbs.AnswerGroupInviteResponse, error)
