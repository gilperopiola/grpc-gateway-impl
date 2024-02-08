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

/* Welcome~! - Here begins this simple Go implementation of the grpc-gateway framework. With gRPC, we design our service in a .proto file and then the server and client code is automatically generated.
 * With grpc-gateway we can expose our gRPC service as a RESTful HTTP API, defining routes and verbs with annotations on the protofile.
 * Then we just generate the gateway code and run it alongside the gRPC server. The gateway will translate HTTP requests to gRPC calls and vice versa, handling input automatically.
 * We also use protovalidate to define input rules on the .proto file itself for each request.

 * Then I wanted to simulate a standard backend architecture and trying to fit everything together nicely.
 */

func main() {

	var (
		protoValidator = initProtoValidator()
		serviceLayer   = initServiceLayer()
		transportLayer = initTransportLayer(protoValidator, serviceLayer)
	)

	grpcServer := initGRPCServer(transportLayer)
	httpGateway := initHTTPGateway()

	go runGRPCServer(grpcServer)
	runHTTPGateway(httpGateway)
}

/* API */

type API struct {
	usersPB.UnimplementedUsersServiceServer
	Transport transport.TransportLayer
}

func (s *API) Signup(ctx context.Context, in *usersPB.SignupRequest) (*usersPB.SignupResponse, error) {
	return s.Transport.Signup(ctx, in)
}

func (s *API) Login(ctx context.Context, in *usersPB.LoginRequest) (*usersPB.LoginResponse, error) {
	return s.Transport.Login(ctx, in)
}

/* Helpers */

func initProtoValidator() *protovalidate.Validator {
	protoValidator, err := protovalidate.New()
	if err != nil {
		log.Fatalf("Failed to initialize validator: %v", err)
	}
	return protoValidator
}

func initServiceLayer() service.ServiceLayer {
	return service.NewService()
}

func initTransportLayer(protoValidator *protovalidate.Validator, serviceLayer service.ServiceLayer) transport.TransportLayer {
	return transport.NewTransport(protoValidator, serviceLayer)
}

func initGRPCServer(transport transport.TransportLayer) *grpc.Server {
	grpcServer := grpc.NewServer()
	usersPB.RegisterUsersServiceServer(grpcServer, &API{Transport: transport})
	return grpcServer
}

func initHTTPGateway() *runtime.ServeMux {
	mux := runtime.NewServeMux()
	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}
	if err := usersPB.RegisterUsersServiceHandlerFromEndpoint(context.Background(), mux, ":50051", opts); err != nil {
		log.Fatalf("Failed to start HTTP gateway: %v", err)
	}
	return mux
}

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

func runHTTPGateway(httpGateway *runtime.ServeMux) {
	log.Println("Running HTTP!")
	if err := http.ListenAndServe(":8080", httpGateway); err != nil {
		log.Fatalf("Failed to serve HTTP: %v", err)
	}
}
