package main

import (
	"context"
	"log"
	"net"
	"net/http"

	usersPB "github.com/gilperopiola/grpc-gateway-impl/pkg/users"
	v1 "github.com/gilperopiola/grpc-gateway-impl/pkg/v1"
	v1Service "github.com/gilperopiola/grpc-gateway-impl/pkg/v1/service"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

/*
	Temp, discovery:

users.pb.gw.go
1. RegisterUsersServiceHandlerClient
2. request_UsersService_Signup_0

users_grpc.pb.go
3. func (c *usersServiceClient) Signup

grpc/server.go
3.2 (s *Server) handleStream
3.5 (s *Server) processUnaryRPC

users_grpc.pb.go
4. _UsersService_Signup_Handler

interceptors.go
5. NewValidationInterceptor (inside func)

IF ERROR UP TO THIS POINT:

	error_handler.go
	6. handleHTTPError (if there is an error before going into the Service layer)

IF NO ERROR UP TO THIS POINT:

	api.go
	6. func (api *API) Signup

	service_users.go
	7. func (s *service) Signup

	users_grpc.pb.go
	8. _UsersService_Signup_Handler (backtracking)

	handler.go
	9. handleForwardResponseOptions

	runtime/mux.go
	10. func (s *ServeMux) ServeHTTP
*/
const (
	gRPCPort = ":50051"
	httpPort = ":8080"
)

/* Welcome~!
 *
 * This is the entrypoint of our app. Here we start the gRPC server and point the HTTP Gateway towards it.
 *
 */

func main() {
	var (
		grpcServer  = initGRPCServer(v1Service.NewService())
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

/* ----------------------------------- */
/*             - gRPC -                */
/* ----------------------------------- */

const (
	errMsgListenGRPC = "Failed to listen gRPC: %v"
	errMsgServeGRPC  = "Failed to serve gRPC: %v"
)

// initGRPCServer initializes the gRPC server and registers the API methods.
// The HTTP Gateway will point towards this server.
func initGRPCServer(service v1Service.ServiceLayer) *grpc.Server {
	var (
		interceptors = v1.GetInterceptors()
		grpcServer   = grpc.NewServer(interceptors)
	)
	usersPB.RegisterUsersServiceServer(grpcServer, &v1.API{Service: service})
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

const (
	errMsgServeHTTP = "Failed to serve HTTP: %v"
	errMsgGateway   = "Failed to start HTTP gateway: %v"
)

// initHTTPGateway initializes the HTTP gateway and registers the API methods there as well.
// The gateway will point towards the gRPC server's port.
func initHTTPGateway() *runtime.ServeMux {
	mux := runtime.NewServeMux(v1.GetHTTPMiddleware()...)
	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	}

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

/* T0D0: Buf file / Dockerfile / Docker-compose / Kubernetes / CI-CD / Tests /
 * Logging / Metrics / Tracing / Security / Caching / Rate limiting / Postman collection
 */
