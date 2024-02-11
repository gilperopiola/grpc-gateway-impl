package main

import (
	"context"
	"log"
	"net"
	"net/http"

	usersPB "github.com/gilperopiola/grpc-gateway-impl/pkg/users"
	v1 "github.com/gilperopiola/grpc-gateway-impl/pkg/v1"
	"github.com/gilperopiola/grpc-gateway-impl/pkg/v1/service"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

/* TODO
 * - README.md v1
 * - Buf file?
 * - Dockerfile
 * - Docker-compose
 * - Kubernetes
 * - CI/CD
 * - Tests
 * - Logging
 * - Metrics
 * - Tracing
 * - Security
 * - Error handling
 * - Versioning
 * - Caching
 * - Rate limiting
 * - Postman collection
 */

/* - Welcome~! - Here begins this simple implementation of the grpc-gateway framework. With gRPC, we design our service in a .proto file and then the server and client code is automatically generated. */

func main() {
	var (
		grpcServer  = initGRPCServer(service.NewService())
		httpGateway = initHTTPGateway()
	)

	go runGRPCServer(grpcServer)
	runHTTPGateway(httpGateway)
}

/* So here we simulate a simple gRPC & HTTP backend API with 2 mock endpoints: Signup and Login.

 * With grpc-gateway we expose our gRPC service as a RESTful HTTP API, defining routes and verbs with annotations on the .proto files.
 * Then we just generate the gateway code and run it alongside the gRPC server. The gateway will translate HTTP requests to gRPC calls, handling input automatically.
 * We also use protovalidate to define input rules on the .proto files themselves for each request, which we enforce using an interceptor.

 * After the validation interceptor runs, requests go through one of the methods on pkg/v1/api.go. From there the Service layer is called, the business logic is executed and the request returns.
 */

const (
	gRPCPort = ":50051"
	httpPort = ":8080"

	errMsgListenGRPC = "Failed to listen gRPC: %v"
	errMsgServeGRPC  = "Failed to serve gRPC: %v"
	errMsgServeHTTP  = "Failed to serve HTTP: %v"
	errMsgGateway    = "Failed to start HTTP gateway: %v"
)

/* ----------------------------------- */
/*             - gRPC -                */
/* ----------------------------------- */

// initGRPCServer initializes the gRPC server and registers the API methods.
// The HTTP Gateway will point towards this server.
func initGRPCServer(serviceLayer service.ServiceLayer) *grpc.Server {
	interceptors := []grpc.ServerOption{v1.NewValidationInterceptor()}
	grpcServer := grpc.NewServer(interceptors...)
	usersPB.RegisterUsersServiceServer(grpcServer, &v1.API{Service: serviceLayer})
	return grpcServer
}

// runGRPCServer runs the gRPC server on a given port.
// It listens for incoming gRPC requests and serves them.
func runGRPCServer(grpcServer *grpc.Server) {
	log.Println("Running gRPC!")

	lis, err := net.Listen("tcp", gRPCPort)
	if err != nil {
		log.Fatalf(errMsgListenGRPC, err)
	}

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf(errMsgServeGRPC, err)
	}
}

/* ----------------------------------- */
/*             - HTTP -                */
/* ----------------------------------- */

// initHTTPGateway initializes the HTTP gateway and registers the API methods there as well.
// The gateway will point towards the gRPC server's port.
func initHTTPGateway() *runtime.ServeMux {
	mux := runtime.NewServeMux()
	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}

	if err := usersPB.RegisterUsersServiceHandlerFromEndpoint(context.Background(), mux, gRPCPort, opts); err != nil {
		log.Fatalf(errMsgGateway, err)
	}

	return mux
}

// runHTTPGateway runs the HTTP gateway on a given port.
// It listens for incoming HTTP requests and serves them.
func runHTTPGateway(httpGateway *runtime.ServeMux) {
	log.Println("Running HTTP!")

	if err := http.ListenAndServe(httpPort, httpGateway); err != nil {
		log.Fatalf(errMsgServeHTTP, err)
	}
}
