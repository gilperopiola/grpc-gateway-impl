package middleware

import (
	"context"
	"log"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
)

/* ----------------------------------- */
/*         - HTTP Middleware -         */
/* ----------------------------------- */

// NewErrorHandlerMiddleware returns a ServeMuxOption that sets a custom error handler for the HTTP gateway.
func NewHTTPErrorHandler() runtime.ServeMuxOption {
	return runtime.WithErrorHandler(handleHTTPError)
}

// NewHTTPLogger returns a ServeMuxOption that logs every HTTP request.
func NewHTTPLogger() runtime.ServeMuxOption {
	return runtime.WithForwardResponseOption(logHTTPResponse)
}

// NewHTTPResponseModifier returns a ServeMuxOption that modifies the HTTP response before it's written to the client.
func NewHTTPResponseModifier() runtime.ServeMuxOption {
	return runtime.WithForwardResponseOption(modifyHTTPResponseHeaders)
}

/* ----------------------------------- */
/*         - Implementations -         */
/* ----------------------------------- */

// logHTTPResponse logs every HTTP request.
func logHTTPResponse(ctx context.Context, rw http.ResponseWriter, msg proto.Message) error {
	log.Printf("HTTP request: %s %s\n", "POST", "/asd")
	return nil
}

// modifyHTTPResponseHeaders executes before the response is written to the client.
func modifyHTTPResponseHeaders(ctx context.Context, rw http.ResponseWriter, resp protoreflect.ProtoMessage) error {

	// Delete gRPC-related headers:
	rw.Header().Del("Grpc-Metadata-Content-Type")

	// Add security-related headers:
	rw.Header().Set("Content-Security-Policy", "default-src 'self'")
	rw.Header().Set("X-XSS-Protection", "1; mode=block")
	rw.Header().Set("X-Frame-Options", "SAMEORIGIN")
	rw.Header().Set("X-Content-Type-Options", "nosniff")
	rw.Header().Set("Strict-Transport-Security", "max-age=31536000; includeSubDomains; preload")

	return nil
}
