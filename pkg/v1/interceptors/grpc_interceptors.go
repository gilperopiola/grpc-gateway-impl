package interceptors

import (
	"github.com/bufbuild/protovalidate-go"
	"google.golang.org/grpc"
)

/* ----------------------------------- */
/*        - gRPC Interceptors -        */
/* ----------------------------------- */

// GetInterceptorsAsServerOption returns a gRPC server option that chains all interceptors together.
// These may be gRPC interceptors, but they are also executed through HTTP calls.
func GetInterceptorsAsServerOption(protoValidator *protovalidate.Validator) grpc.ServerOption {
	return grpc.ChainUnaryInterceptor(
		newValidationInterceptor(protoValidator),
	)
}
