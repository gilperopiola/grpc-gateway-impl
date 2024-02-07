package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"

	"github.com/bufbuild/protovalidate-go"
	usersPB "github.com/gilperopiola/grpc-gateway-impl/pkg/users"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type API struct {
	usersPB.UnimplementedUsersServiceServer

	Transport TransportLayer
}

func initProtoValidator() *protovalidate.Validator {
	protoValidator, err := protovalidate.New()
	if err != nil {
		log.Fatalf("Failed to initialize validator:", err)
	}
	return protoValidator
}

func initTransportLayer(protoValidator *protovalidate.Validator) TransportLayer {
	return &transport{
		protoValidator: protoValidator,
		service:        &service{},
	}
}

func initGRPCServer(protoValidator *protovalidate.Validator) *grpc.Server {
	grpcServer := grpc.NewServer()
	usersPB.RegisterUsersServiceServer(grpcServer, &API{
		Transport: initTransportLayer(protoValidator),
	})
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

func main() {

	// Init dependencies (validator)
	protoValidator := initProtoValidator()
	grpcServer := initGRPCServer(protoValidator)
	mux := initHTTPGateway()

	// Start gRPC server
	go func() {
		lis, err := net.Listen("tcp", ":50051")
		if err != nil {
			log.Fatalf("Failed to listen gRPC: %v", err)
		}
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatalf("Failed to serve gRPC: %v", err)
		}
	}()

	// Start HTTP server (grpc-gateway)
	log.Println("Running HTTP!")
	http.ListenAndServe(":8080", mux)
}

func (s *API) Signup(ctx context.Context, in *usersPB.SignupRequest) (*usersPB.SignupResponse, error) {
	return s.Transport.Signup(ctx, in)
}

/* pkg/entities */

type SignupRequest struct {
	Email    string
	Password string
}

type SignupResponse struct {
	Email string
}

/* pkg/service */

type ServiceLayer interface {
	Signup(ctx context.Context, in SignupRequest) (SignupResponse, error)
}

type service struct {
}

/* pkg/transport */

type TransportLayer interface {
	Signup(ctx context.Context, in *usersPB.SignupRequest) (*usersPB.SignupResponse, error)
}

type transport struct {
	protoValidator *protovalidate.Validator

	service     ServiceLayer
	toService   Transport2ServiceAdapter
	toTransport Service2TransportAdapter
}

func (t *transport) Signup(ctx context.Context, transportRequest *usersPB.SignupRequest) (*usersPB.SignupResponse, error) {
	if err := t.protoValidator.Validate(transportRequest); err != nil {
		fmt.Println("validation failed:", err)
	}

	serviceRequest, err := t.toService.Signup(transportRequest)
	if err != nil {
		return nil, err
	}

	serviceResponse, err := t.service.Signup(ctx, serviceRequest)
	if err != nil {
		return nil, err
	}

	transportResponse, err := t.toTransport.Signup(serviceResponse)
	if err != nil {
		return nil, err
	}

	return transportResponse, nil
}

/* pkg/adapters */

type Transport2ServiceAdapter interface {
}

type Service2TransportAdapter interface {
}
