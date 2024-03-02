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
	"github.com/gilperopiola/grpc-gateway-impl/pkg/v1/service"
	"github.com/gilperopiola/grpc-gateway-impl/server/config"
	"github.com/gilperopiola/grpc-gateway-impl/server/security"

	"github.com/bufbuild/protovalidate-go"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

/* ----------------------------------- */
/*         - Main Application -        */
/* ----------------------------------- */

// App holds every dependency we need to init and run the servers, and the servers themselves.
// It has an embedded *v1.API, which is our concrete implementation of the gRPC API.
type App struct {
	*v1.API

	// Config holds the configuration of our API.
	// Service holds the business logic of our API.
	Config  *config.Config
	Service service.Service

	// GRPCServer and HTTPGateway are the servers we are going to run.
	GRPCServer  *grpc.Server
	HTTPGateway *http.Server

	// GRPCInterceptors run before or after gRPC calls.
	// Right now we only have a Logger and a Validator.
	//
	// GRPCDialOptions configure the communication between the HTTP Gateway and the gRPC Server.
	// Needed to establish a secure connection.
	GRPCInterceptors []grpc.ServerOption
	GRPCDialOptions  []grpc.DialOption

	// HTTPMiddleware are the ServeMuxOptions that run before or after HTTP calls.
	// Right now we only have an Error Handler and a Response Modifier.
	//
	// HTTPMiddlewareWrapper are the middleware that wrap around the HTTP server.
	// Right now we only have a Logger.
	//
	// They are divided into two different types because the ServeMuxOptions are used to configure the ServeMux,
	// and the Wrapper is used to wrap the ServeMux with middleware.
	HTTPMiddleware        []runtime.ServeMuxOption
	HTTPMiddlewareWrapper middleware.MuxWrapperFunc

	// Logger is used to log every gRPC request that comes in through the gRPC
	// It's used on an interceptor.
	Logger *zap.Logger

	// ProtoValidator is used to validate the incoming gRPC & HTTP requests.
	// It uses the bufbuild/protovalidate library to enforce the validation rules written in the .proto files.
	ProtoValidator *protovalidate.Validator

	// TLSServerCertificate is a pool of certificates to use for the
	// It is used to validate the gRPC server's certificate on the HTTP Gateway calls.
	// TLSServerCredentials and TLSClientCredentials are used to establish the
	// secure connection between the HTTP Gateway and the gRPC Server.
	TLSServerCertificate *x509.CertPool
	TLSServerCredentials credentials.TransportCredentials
	TLSClientCredentials credentials.TransportCredentials
}

// NewApp returns a new App with the given configuration.
func NewApp(config *config.Config) *App {
	return &App{Config: config}
}

// InitGeneralDependencies initializes stuff.
// Logger, Validator and TLS.
func (a *App) InitGeneralDependencies() {
	a.Logger = v1.NewLogger(a.Config.IsProd, v1.NewLoggerOptions())
	a.ProtoValidator = interceptors.NewProtoValidator()

	tlsConfig := a.Config.TLS
	{
		// TLS Server objects are only loaded if TLS is enabled.
		if tlsConfig.Enabled {
			a.TLSServerCertificate = security.NewTLSCertPool(tlsConfig.CertPath)
			a.TLSServerCredentials = security.NewServerTransportCredentials(tlsConfig.CertPath, tlsConfig.KeyPath)
		}

		// TLS Client Credentials are always loaded, with an insecure option if TLS is not enabled.
		a.TLSClientCredentials = security.NewClientTransportCredentials(tlsConfig.Enabled, a.TLSServerCertificate)
	}
}

// InitGRPCAndHTTPDependencies initializes gRPC and HTTP stuff.
func (a *App) InitGRPCAndHTTPDependencies() {

	// gRPC Interceptors and Dial Options.
	a.GRPCInterceptors = interceptors.GetAll(a.Config, a.Logger, a.ProtoValidator, a.TLSServerCredentials)
	a.GRPCDialOptions = getAllDialOptions(a.TLSClientCredentials)

	// HTTP Middleware and Mux Wrapper.
	a.HTTPMiddleware = middleware.GetAll()
	a.HTTPMiddlewareWrapper = middleware.GetAllWrapped(a.Logger)
}

// InitAPIAndServers initializes the API, the Service and the Servers.
func (a *App) InitAPIAndServers() {

	// Service and API.
	a.Service = service.NewService()
	a.API = v1.NewAPI(a.Service)

	// gRPC and HTTP Servers.
	a.GRPCServer = initGRPCServer(a.API, a.GRPCInterceptors)
	a.HTTPGateway = initHTTPGateway(a.Config.MainConfig, a.HTTPMiddleware, a.HTTPMiddlewareWrapper, a.GRPCDialOptions)
}

// Run runs the gRPC and HTTP servers.
func (a *App) Run() {
	runGRPCServer(a.GRPCServer, a.Config.GRPCPort)
	runHTTPGateway(a.HTTPGateway)
}

// WaitForGracefulShutdown waits for a SIGINT or SIGTERM to gracefully shutdown the servers.
func (a *App) WaitForGracefulShutdown() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM) // SIGINT and SIGTERM
	<-c

	shutdownGRPCServer(a.GRPCServer)
	shutdownHTTPGateway(a.HTTPGateway)

	log.Println("Servers stopped! Bye bye~")
}
