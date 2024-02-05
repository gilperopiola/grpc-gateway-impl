package main

import (
	"context"
	"log"
	"net"
	"net/http"

	usersPB "github.com/gilperopiola/grpc-gateway-impl/pkg/users"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
)

type server struct {
	usersPB.UnimplementedUsersServiceServer
}

func (s *server) Signup(ctx context.Context, in *usersPB.SignupRequest) (*usersPB.SignupResponse, error) {
	return &usersPB.SignupResponse{Message: "Hello " + in.Username}, nil
}
func (s *server) Login(ctx context.Context, in *usersPB.LoginRequest) (*usersPB.LoginResponse, error) {
	return &usersPB.LoginResponse{Message: "Hello " + in.Username}, nil
}

func main() {
	grpcServer := grpc.NewServer()
	usersPB.RegisterUsersServiceServer(grpcServer, &server{})

	// Start gRPC server
	go func() {
		lis, err := net.Listen("tcp", ":50051")
		if err != nil {
			log.Fatalf("failed to listen: %v", err)
		}
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatalf("failed to serve: %v", err)
		}
	}()

	// Start HTTP server (grpc-gateway)
	mux := runtime.NewServeMux()
	opts := []grpc.DialOption{grpc.WithInsecure()}
	err := usersPB.RegisterUsersServiceHandlerFromEndpoint(context.Background(), mux, ":50051", opts)
	if err != nil {
		log.Fatalf("failed to start HTTP gateway: %v", err)
	}
	log.Println("Running HTTP!")
	http.ListenAndServe(":8080", mux)
}
