package middleware

import (
	"context"
	"net/http"

	v1 "github.com/gilperopiola/grpc-gateway-impl/pkg/v1"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"go.uber.org/zap"
	"google.golang.org/protobuf/reflect/protoreflect"
)

// MuxWrapperFunc is a middleware that wraps around the HTTP Server's mux.
type MuxWrapperFunc func(next http.Handler) http.Handler

/* ----------------------------------- */
/*         - HTTP Middleware -         */
/* ----------------------------------- */
// Some middleware are passed as ServeMuxOptions when the mux is created,
// and some are wrapped around the mux after its creation.

// GetAll returns all the HTTP middleware that are used as ServeMuxOptions.
func GetAll() []runtime.ServeMuxOption {
	return []runtime.ServeMuxOption{
		runtime.WithForwardResponseOption(setHTTPResponseHeaders),
		runtime.WithErrorHandler(handleHTTPError),
	}
}

// GetAllWrapped returns the middleware to be wrapped around the HTTP Server when it's created.
// It handles CORS and logs the HTTP Request's info when it finishes executing.
// It's used to wrap the mux with middleware.
func GetAllWrapped(logger *zap.Logger) MuxWrapperFunc {
	return func(next http.Handler) http.Handler {
		return handleCORS(
			v1.LogHTTP(next, logger),
		)
	}
}

// handleCORS adds CORS headers to the response and handles preflight requests.
func handleCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")

		// Handle preflight requests.
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		// Pass down the chain to next handler if not OPTIONS.
		next.ServeHTTP(w, r)
	})
}

// setHTTPResponseHeaders executes before the response is written to the client.
func setHTTPResponseHeaders(_ context.Context, rw http.ResponseWriter, _ protoreflect.ProtoMessage) error {
	for _, headerToBeDeleted := range httpResponseHeadersToDelete {
		rw.Header().Del(headerToBeDeleted)
	}
	for headerKey, headerValue := range httpResponseHeadersToAdd {
		rw.Header().Set(headerKey, headerValue)
	}
	return nil
}

var httpResponseHeadersToAdd = map[string]string{
	"Content-Security-Policy":   "default-src 'self'",
	"Content-Type":              "application/json",
	"Strict-Transport-Security": "max-age=31536000; includeSubDomains; preload",
	"X-Content-Type-Options":    "nosniff",
	"X-Frame-Options":           "SAMEORIGIN",
	"X-XSS-Protection":          "1; mode=block",
}
var httpResponseHeadersToDelete = []string{"Grpc-Metadata-Content-Type"}
