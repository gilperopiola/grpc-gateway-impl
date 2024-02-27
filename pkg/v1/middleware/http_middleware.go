package middleware

import (
	"context"
	"net/http"
	"time"

	v1 "github.com/gilperopiola/grpc-gateway-impl/pkg/v1"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"go.uber.org/zap"
	"google.golang.org/protobuf/reflect/protoreflect"
)

/* ----------------------------------- */
/*         - HTTP Middleware -         */
/* ----------------------------------- */

func GetAll() v1.MiddlewareI {
	return v1.MiddlewareI{
		runtime.WithErrorHandler(handleHTTPError),
		runtime.WithForwardResponseOption(modifyHTTPResponseHeaders),
	}
}

// LogHTTP doesn't work like a middleware (it's wrapped around the HTTP server when it's created),
// but I think this belongs here.
func LogHTTP(logger *zap.Logger) func(next http.Handler) http.Handler {
	sugar := logger.Sugar()

	return func(next http.Handler) http.Handler {
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
