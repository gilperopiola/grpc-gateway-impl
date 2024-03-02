package middleware

import (
	"context"
	"net/http"

	v1 "github.com/gilperopiola/grpc-gateway-impl/pkg/v1"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"go.uber.org/zap"
	"google.golang.org/protobuf/reflect/protoreflect"
)

/* ----------------------------------- */
/*         - HTTP Middleware -         */
/* ----------------------------------- */
// Some middleware are passed as ServeMuxOptions when the mux is created,
// and some are wrapped around the mux after its creation.

// GetAll returns all the HTTP middleware that are used as ServeMuxOptions.
func GetAll() []runtime.ServeMuxOption {
	return []runtime.ServeMuxOption{
		runtime.WithErrorHandler(handleHTTPError),
		runtime.WithForwardResponseOption(setHTTPResponseHeaders),
	}
}

// MuxWrapperFunc is a middleware that wraps around the HTTP Server's mux.
type MuxWrapperFunc func(next http.Handler) http.Handler

// GetMuxWrapperFunc is wrapped around the HTTP server when it's created
// and logs the HTTP Request's info when it finishes executing.
// It's used to wrap the mux with middleware.
func GetMuxWrapperFunc(logger *zap.Logger) MuxWrapperFunc {
	return func(next http.Handler) http.Handler {
		return v1.LogHTTP(next, logger)
	}
}

// setHTTPResponseHeaders executes before the response is written to the client.
func setHTTPResponseHeaders(ctx context.Context, rw http.ResponseWriter, resp protoreflect.ProtoMessage) error {
	for _, headerToBeDeleted := range defaultResponseHeadersToBeDeleted {
		rw.Header().Del(headerToBeDeleted)
	}

	for headerKey, headerValue := range defaultResponseHeaders {
		rw.Header().Set(headerKey, headerValue)
	}

	return nil
}

var defaultResponseHeaders = map[string]string{
	"Content-Security-Policy":   "default-src 'self'",
	"X-XSS-Protection":          "1; mode=block",
	"X-Frame-Options":           "SAMEORIGIN",
	"X-Content-Type-Options":    "nosniff",
	"Strict-Transport-Security": "max-age=31536000; includeSubDomains; preload",
}

var defaultResponseHeadersToBeDeleted = []string{
	"Grpc-Metadata-Content-Type",
}
