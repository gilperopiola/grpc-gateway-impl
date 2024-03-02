package server

import (
	"testing"

	"github.com/gilperopiola/grpc-gateway-impl/server/config"
)

func TestInitGeneralDependencies(t *testing.T) {
	app := NewApp(config.LoadConfig())
	app.InitGeneralDependencies()

	if app.Logger == nil {
		t.Errorf("Expected Logger to be initialized, got nil")
	}
	if app.ProtoValidator == nil {
		t.Errorf("Expected ProtoValidator to be initialized, got nil")
	}
}
