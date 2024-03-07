package server

import (
	v1 "github.com/gilperopiola/grpc-gateway-impl/pkg/v1"
	"github.com/gilperopiola/grpc-gateway-impl/pkg/v1/interceptors"
	"github.com/gilperopiola/grpc-gateway-impl/pkg/v1/middleware"
	"github.com/gilperopiola/grpc-gateway-impl/pkg/v1/service"
	"github.com/gilperopiola/grpc-gateway-impl/server/config"
)

func (a *App) NewConfig() {
	a.Cfg = config.New()
}

func (a *App) NewGRPCServer() {
	a.GRPCServer = newGRPCServer(a.Cfg.GRPCPort, a.API, a.GRPCInterceptors)
	a.GRPCServer.Init()
}

func (a *App) NewHTTPGateway() {
	a.HTTPGateway = newHTTPGateway(a.Cfg.MainConfig, a.HTTPMiddleware, a.HTTPMiddlewareWrapper, a.GRPCDialOptions)
	a.HTTPGateway.Init()
}

func (a *App) NewAPI() {
	a.API = v1.NewAPI(a.Service)
}

func (a *App) NewService() {
	a.Service = service.NewService()
}

func (a *App) NewGRPCInterceptors() {
	a.GRPCInterceptors = interceptors.All(a.Cfg, a.Logger, a.ProtoValidator, a.TLSServerCreds)
}

func (a *App) NewGRPCDialOptions() {
	a.GRPCDialOptions = AllDialOptions(a.TLSClientCreds)
}

func (a *App) NewHTTPMiddleware() {
	a.HTTPMiddleware = middleware.All()
}

func (a *App) NewHTTPMiddlewareWrapper() {
	a.HTTPMiddlewareWrapper = middleware.Wrapper(a.Logger)
}

func (a *App) NewProtoValidator() {
	a.ProtoValidator = interceptors.NewProtoValidator()
}

func (a *App) NewLogger() {
	a.Logger = v1.NewLogger(a.Cfg.IsProd, a.LoggerOptions)
}

func (a *App) NewLoggerOptions() {
	a.LoggerOptions = v1.NewLoggerOptions()
}

func (a *App) NewTLSServerCert() {
	a.TLSServerCert = config.NewTLSCertPool(a.Cfg.TLS.CertPath)
}

func (a *App) NewTLSServerCreds() {
	a.TLSServerCreds = config.NewServerTransportCredentials(a.Cfg.TLS.CertPath, a.Cfg.TLS.KeyPath)
}

func (a *App) NewTLSClientCreds() {
	a.TLSClientCreds = config.NewClientTransportCredentials(a.Cfg.TLS.Enabled, a.TLSServerCert)
}
