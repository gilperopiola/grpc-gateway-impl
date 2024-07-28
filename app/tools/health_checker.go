package tools

import (
	"context"
	"errors"

	"github.com/gilperopiola/grpc-gateway-impl/app/core"
)

var _ core.HealthChecker = (*healthChecker)(nil)

type healthChecker struct {
	healthCheckFn ServiceFunc
}

func NewHealthChecker(serviceFn ServiceFunc) *healthChecker {
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
