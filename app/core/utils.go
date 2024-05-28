package core

import (
	"context"
	"net/http"
	"strconv"
	"time"

	"github.com/gilperopiola/grpc-gateway-impl/app/core/pbs"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

// -> Yes, we have a utils file. List of reasons why we shouldn't have one:
// -

// ...
// Get it? It's empty. Utils are fine.

type (
	Ctx = context.Context

	AuthSvc   = pbs.AuthServiceServer
	UsersSvc  = pbs.UsersServiceServer
	GroupsSvc = pbs.GroupsServiceServer

	GRPCInfo             = grpc.UnaryServerInfo
	GRPCHandler          = grpc.UnaryHandler
	GRPCInterceptors     = []grpc.UnaryServerInterceptor
	GRPCServerOptions    = []grpc.ServerOption
	GRPCDialOptions      = []grpc.DialOption
	GRPCServiceRegistrar = grpc.ServiceRegistrar

	HTTPMultiplexer = runtime.ServeMux

	TLSCredentials = credentials.TransportCredentials
)

/* -~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~- */

func NewCtx() Ctx {
	return context.Background()
}

func NewCtxWithTimeout(duration time.Duration) (Ctx, context.CancelFunc) {
	return context.WithTimeout(NewCtx(), duration)
}

func IfErrThen(err error, do func()) {
	if err != nil {
		do()
	}
}

func ToIntAndErr(s string, err error) (int, error) {
	if err != nil {
		return 0, err
	}
	return strconv.Atoi(s)
}

/* -~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~- */

var _ http.ResponseWriter = (*CustomResponseWriter)(nil)

type CustomResponseWriter struct {
	http.ResponseWriter
	Status int
}

func (crw *CustomResponseWriter) WriteHeader(statusCode int) {
	crw.Status = statusCode
	crw.ResponseWriter.WriteHeader(statusCode)
}

/* -~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~- */

type Int32Slice []int32

func (int32Slice Int32Slice) ToIntSlice() []int {
	intSlice := make([]int, len(int32Slice))
	for i, int32Value := range int32Slice {
		intSlice[i] = int(int32Value)
	}
	return intSlice
}


