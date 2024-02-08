package main

import (
	"context"
	"log"
	"net"
	"net/http"

	"github.com/gilperopiola/grpc-gateway-impl/pkg/service"
	"github.com/gilperopiola/grpc-gateway-impl/pkg/transport"
	usersPB "github.com/gilperopiola/grpc-gateway-impl/pkg/users"

	"github.com/bufbuild/protovalidate-go"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

/* - Welcome~! - Here begins this simple implementation of the grpc-gateway framework. With gRPC, we design our service in a .proto file and then the server and client code is automatically generated. */

func main() {
	var (
		transportLayer = initTransportLayer(initProtoValidator(), initServiceLayer())
		grpcServer     = initGRPCServer(transportLayer)
		httpGateway    = initHTTPGateway()
	)

	go runGRPCServer(grpcServer)
	runHTTPGateway(httpGateway)
}

/* With grpc-gateway we can expose our gRPC service as a RESTful HTTP API, defining routes and verbs with annotations on the protofile.
 * Then we just generate the gateway code and run it alongside the gRPC server. The gateway will translate HTTP requests to gRPC calls and vice versa, handling input automatically.
 * We also use protovalidate to define input rules on the .proto file itself for each request, which is really useful.

 * So here we simulate a simple gRPC & HTTP backend API with 2 endpoints: Signup and Login. They don't really do anything, we just showcase the architecture.
 *  - First, we define the .proto file with the service and messages and automatically generate the rest of code.
 *  - Then we start the application, initializing the .proto validator, the service layer and the transport layer.
 *  - We run the gRPC server on port 50051 and then the HTTP gateway pointing towards it on 8080, they run in parallel.
 */

// runGRPCServer runs the gRPC server on port 50051.
// It listens for incoming gRPC requests and serves them.
func runGRPCServer(grpcServer *grpc.Server) {
	log.Println("Running gRPC!")
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("Failed to listen gRPC: %v", err)
	}
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve gRPC: %v", err)
	}
}

// runHTTPGateway runs the HTTP gateway on port 8080.
// It listens for incoming HTTP requests and serves them.
func runHTTPGateway(httpGateway *runtime.ServeMux) {
	log.Println("Running HTTP!")
	if err := http.ListenAndServe(":8080", httpGateway); err != nil {
		log.Fatalf("Failed to serve HTTP: %v", err)
	}
}

/* Requests enter through the API methods down below, which are basically the gRPC methods. From there the Transport layer is called, validating the input and using a generic function to handle
 * all types of requests in the same way. Transport uses an adapter to convert the request from the PB format to our custom entities format, and then calls the Service layer with that request.
 * The call to the Service layer returns a custom entities format response, which the Transport converts back to the PB format and then the API returns it.
 */

// API is our concrete implementation of the gRPC API defined in the .proto files.
// It implements a handler for each API method, connecting it with the Transport layer.
type API struct {
	usersPB.UnimplementedUsersServiceServer
	Transport transport.TransportLayer
}

// Signup is the handler for the Signup API method. Both gRPC and HTTP calls will trigger this method.
func (s *API) Signup(ctx context.Context, in *usersPB.SignupRequest) (*usersPB.SignupResponse, error) {
	return s.Transport.Signup(ctx, in)
}

// Login is the handler for the Login API method. Both gRPC and HTTP calls will trigger this method.
func (s *API) Login(ctx context.Context, in *usersPB.LoginRequest) (*usersPB.LoginResponse, error) {
	return s.Transport.Login(ctx, in)
}

/* ----------------------------------- */
/*    - Initialization Functions -     */
/* ----------------------------------- */

// initProtoValidator initializes the .proto validator.
// With it we can define input rules on the .proto file itself for each request to follow.
func initProtoValidator() *protovalidate.Validator {
	protoValidator, err := protovalidate.New()
	if err != nil {
		log.Fatalf("Failed to initialize validator: %v", err)
	}
	return protoValidator
}

// initServiceLayer initializes the service layer.
// There we define the business logic of our application, sometimes incorporating a RepositoryLayer on there as well.
func initServiceLayer() service.ServiceLayer {
	return service.NewService()
}

// initTransportLayer initializes the transport layer.
// There we define the API methods, which basically are the gRPC methods that HTTP calls will also trigger.
// It converts requests and responses from the PB format to our custom entities format and vice versa, using the entities format to call the Service layer.
func initTransportLayer(protoValidator *protovalidate.Validator, serviceLayer service.ServiceLayer) transport.TransportLayer {
	return transport.NewTransport(protoValidator, serviceLayer)
}

// initGRPCServer initializes the gRPC server and registers the API methods.
// The HTTP Gateway will point towards this server.
func initGRPCServer(transport transport.TransportLayer) *grpc.Server {
	grpcServer := grpc.NewServer()
	usersPB.RegisterUsersServiceServer(grpcServer, &API{Transport: transport})
	return grpcServer
}

// initHTTPGateway initializes the HTTP gateway and registers the API methods there as well.
// The gateway will point towards the gRPC server, on port 50051.
func initHTTPGateway() *runtime.ServeMux {
	mux := runtime.NewServeMux()
	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}
	if err := usersPB.RegisterUsersServiceHandlerFromEndpoint(context.Background(), mux, ":50051", opts); err != nil {
		log.Fatalf("Failed to start HTTP gateway: %v", err)
	}
	return mux
}
