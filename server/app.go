package server

import (
	"crypto/x509"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	v1 "github.com/gilperopiola/grpc-gateway-impl/pkg/v1"
	"github.com/gilperopiola/grpc-gateway-impl/pkg/v1/interceptors"
	"github.com/gilperopiola/grpc-gateway-impl/pkg/v1/middleware"

	"github.com/bufbuild/protovalidate-go"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

// App holds every dependency we need to init and run the servers, and the servers themselves.
// It has an embedded *v1.API, which is our concrete implementation of the gRPC API.
type App struct {
	*v1.API

	// Config holds the configuration of our API.
	// Service holds the business logic of our API.
	// GRPCServer and HTTPGateway are the servers we are going to run.
	Config      *v1.APIConfig
	Service     v1.ServiceLayer
	GRPCServer  *grpc.Server
	HTTPGateway *http.Server

	// GRPCInterceptors run before or after gRPC calls.
	// Right now we only have a Logger and a Validator.
	//
	// GRPCDialOptions configure the communication between the HTTP Gateway and the gRPC server.
	// Needed to establish a secure connection to the gRPC server.
	//
	// HTTPMiddleware run before or after HTTP calls.
	// Right now we only have an Error Handler and a Response Modifier.
	GRPCInterceptors grpc.ServerOption
	GRPCDialOptions  v1.GRPCDialOptionsI
	HTTPMiddleware   v1.MiddlewareI

	// Logger is used to log every gRPC request that comes in through the gRPC server.
	// It's used on an interceptor.
	//
	// ProtoValidator is used to validate the incoming gRPC & HTTP requests.
	// It uses the bufbuild/protovalidate library to enforce the validation rules written in the .proto files.
	//
	// TLSCertPool is a pool of certificates to use for the server.
	// It is used to validate the gRPC server's certificate on the HTTP Gateway calls.
	Logger         *zap.Logger
	ProtoValidator *protovalidate.Validator
	TLSCertPool    *x509.CertPool
}

// Prepare initializes the App and its dependencies.
func (app *App) Prepare() {
	app.Logger = newLogger(app.Config.IsProd)
	app.ProtoValidator = v1.NewProtoValidator()
	app.TLSCertPool = LoadTLSCertPool(app.Config.TLS.CertPath)

	app.HTTPMiddleware = middleware.GetAll()
	app.GRPCDialOptions = interceptors.GetAllDialOptions(app.Config.TLS, app.TLSCertPool)
	app.GRPCInterceptors = interceptors.GetAll(app.Logger, app.ProtoValidator)

	app.Service = v1.NewService()
	app.API = v1.NewAPI(app.Service)
	app.GRPCServer = newGRPCServer(app.API, app.Config.TLS, app.GRPCInterceptors)
	app.HTTPGateway = newHTTPGateway(app.Config, app.HTTPMiddleware, app.GRPCDialOptions, app.Logger)
}

func newGRPCServer(api *v1.API, tlsConfig v1.TLSConfig, interceptors grpc.ServerOption) *grpc.Server {
	return InitGRPCServer(api, tlsConfig, interceptors)
}

func newHTTPGateway(config *v1.APIConfig, httpMiddleware v1.MiddlewareI, grpcDialOptions v1.GRPCDialOptionsI, logger *zap.Logger) *http.Server {
	return InitHTTPGateway(config.GRPCPort, config.HTTPPort, httpMiddleware, grpcDialOptions, middleware.LogHTTP(logger))
}

func newLogger(isProd bool) *zap.Logger {
	loggerOptions := []zap.Option{
		zap.AddStacktrace(zap.DPanicLevel), // Add stack trace to panic logs.
	}
	return v1.NewLogger(isProd, loggerOptions)
}

// WaitForGracefulShutdown waits for a SIGINT or SIGTERM to gracefully shutdown the servers.
func (aw *App) WaitForGracefulShutdown() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM) // SIGINT and SIGTERM
	<-c

	ShutdownGRPCServer(aw.GRPCServer)
	ShutdownHTTPServer(aw.HTTPGateway)

	log.Println("Servers stopped! Bye bye~")
}
