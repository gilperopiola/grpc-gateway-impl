package middleware

import (
	"context"
	"net/http"
	"time"

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

// MuxWrapperFn is a middleware that wraps around the HTTP Server's mux.
type MuxWrapperFn func(next http.Handler) http.Handler

// GetMuxWrapperFn is wrapped around the HTTP server when it's created
// and logs the HTTP Request's info when it finishes executing.
// It's used to wrap the mux with middleware.
func GetMuxWrapperFn(logger *zap.Logger) MuxWrapperFn {
	sugar := logger.Sugar()
	return func(next http.Handler) http.Handler {
		return logHTTP(next, sugar)
	}
}

// logHTTP logs the HTTP Request's info when it finishes executing.
func logHTTP(next http.Handler, sugar *zap.SugaredLogger) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)

		sugar.Infow("HTTP Request",
			"path", r.URL.Path,
			"method", r.Method,
			"duration", time.Since(start),
		)
	})
}

// setHTTPResponseHeaders executes before the response is written to the client.
func setHTTPResponseHeaders(ctx context.Context, rw http.ResponseWriter, resp protoreflect.ProtoMessage) error {

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
