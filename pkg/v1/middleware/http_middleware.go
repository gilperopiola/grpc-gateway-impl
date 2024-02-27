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

// GetAll returns all the HTTP middleware that are used as ServeMuxOptions.
func GetAll() []runtime.ServeMuxOption {
	return []runtime.ServeMuxOption{
		runtime.WithErrorHandler(handleHTTPError),
		runtime.WithForwardResponseOption(modifyHTTPResponseHeaders),
	}
}

// GetMuxWrapperFn is wrapped around the HTTP server when it's created
// and logs the HTTP Request's info when it finishes executing.
// It's used to wrap the ServeMux with middleware.
func GetMuxWrapperFn(logger *zap.Logger) muxWrapperFn {
	sugar := logger.Sugar()
	return func(next http.Handler) http.Handler {
		return logHTTP(next, sugar)
	}
}

// muxWrapperFn is a middleware that wraps around the HTTP Server's mux.
type muxWrapperFn func(next http.Handler) http.Handler

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
