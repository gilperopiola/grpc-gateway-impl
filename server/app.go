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

	"github.com/bufbuild/protovalidate-go"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
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
	Config      *config.Config
	Service     service.Service
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
	HTTPMiddlewareWrapper func(http.Handler) http.Handler

	// Logger is used to log every gRPC request that comes in through the gRPC
	// It's used on an interceptor.
	//
	// ProtoValidator is used to validate the incoming gRPC & HTTP requests.
	// It uses the bufbuild/protovalidate library to enforce the validation rules written in the .proto files.
	//
	// ServerTLSCert is a pool of certificates to use for the
	// It is used to validate the gRPC server's certificate on the HTTP Gateway calls.
	Logger         *zap.Logger
	ProtoValidator *protovalidate.Validator
	ServerTLSCert  *x509.CertPool
}

// NewApp returns a new App with the given configuration.
func NewApp(config *config.Config) *App {
	return &App{Config: config}
}

// InitGeneralDependencies initializes stuff.
// Logger, Validator and Server TLS Certificate.
func (a *App) InitGeneralDependencies() {
	a.Logger = newLogger(a.Config.IsProd, newLoggerOptions())
	a.ProtoValidator = interceptors.NewProtoValidator()
	a.ServerTLSCert = newTLSCertPool(a.Config.TLS.CertPath)
}

// InitGRPCAndHTTPDependencies initializes gRPC and HTTP stuff.
func (a *App) InitGRPCAndHTTPDependencies() {
	tlsEnabled, certPath, keyPath := a.Config.TLS.Enabled, a.Config.TLS.CertPath, a.Config.TLS.KeyPath

	// gRPC Interceptors and Dial Options.
	a.GRPCInterceptors = interceptors.GetAll(a.Logger, a.ProtoValidator, tlsEnabled, certPath, keyPath)
	a.GRPCDialOptions = getAllDialOptions(tlsEnabled, a.ServerTLSCert)

	// HTTP Middleware and Mux Wrapper.
	a.HTTPMiddleware = middleware.GetAll()
	a.HTTPMiddlewareWrapper = middleware.GetMuxWrapperFn(a.Logger)
}

// InitAPIAndServers initializes the API, the Service and the Servers.
func (a *App) InitAPIAndServers() {

	// Service and API.
	a.Service = service.NewService()
	a.API = v1.NewAPI(a.Service)

	// gRPC and HTTP Servers.
	a.GRPCServer = initGRPCServer(a.API, a.GRPCInterceptors)
	a.HTTPGateway = initHTTPGateway(a.Config.GRPCPort, a.Config.HTTPPort, a.HTTPMiddleware, a.GRPCDialOptions, a.HTTPMiddlewareWrapper)
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

// newTLSCertPool loads the server's certificate from a file and returns a certificate pool.
// It's a SSL/TLS certificate used to secure the communication between the HTTP Gateway and the gRPC server.
// It must be in a .crt format.
//
// To generate a self-signed certificate, you can use the following command:
// openssl req -x509 -newkey rsa:4096 -keyout server.key -out server.crt -days 365 -nodes -subj '/CN=localhost'
// The certificate must be in the root directory of the project.
func newTLSCertPool(tlsCertPath string) *x509.CertPool {

	// Read certificate.
	cert, err := os.ReadFile(tlsCertPath)
	if err != nil {
		log.Fatalf(errMsgReadingTLSCert_Fatal, err)
	}

	// Create certificate pool.
	if out := x509.NewCertPool(); out.AppendCertsFromPEM(cert) {
		return out
	}

	// Error appending certificate.
	log.Fatalf(errMsgAppendingTLSCert_Fatal)
	return nil
}

const (
	errMsgReadingTLSCert_Fatal   = "Failed to read TLS certificate: %v" // Fatal error.
	errMsgAppendingTLSCert_Fatal = "Failed to append TLS certificate"   // Fatal error.
)
