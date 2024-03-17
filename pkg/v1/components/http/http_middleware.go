package http

import (
	"context"
	"net/http"

	"github.com/gilperopiola/grpc-gateway-impl/pkg/v1/components/common"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/protobuf/reflect/protoreflect"
)

/* ----------------------------------- */
/*         - HTTP Middleware -         */
/* ----------------------------------- */

/* Some middleware are passed as ServeMuxOptions when the mux is created.
/* Some are wrapped around the Mux afterwards. */

// MuxWrapperFunc is a middleware that wraps around the HTTP Gateway's mux.
type MuxWrapperFunc func(next http.Handler) http.Handler

// AllMiddleware returns all the HTTP middleware that are used as ServeMuxOptions.
func AllMiddleware() []runtime.ServeMuxOption {
	return []runtime.ServeMuxOption{
		runtime.WithErrorHandler(handleHTTPError), // Stops other middleware if an error happens.
		runtime.WithForwardResponseOption(setHTTPRespHeaders),
	}
}

// MiddlewareWrapper returns the middleware to be wrapped around the HTTP Gateway's Mux.
func MiddlewareWrapper(logger *common.Logger) MuxWrapperFunc {
	return func(next http.Handler) http.Handler {
		return handleCORS(
			logger.LogHTTP(next),
		)
	}
}

// handleCORS adds CORS headers to the response and handles preflight requests.
func handleCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		for key, value := range corsHeadersToAdd { // Add headers.
			w.Header().Set(key, value)
		}

		if r.Method == "OPTIONS" { // Handle preflight requests.
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r) // Pass down the chain to next handler if not OPTIONS.
	})
}

// setHTTPRespHeaders executes before the response is written to the client.
// It's also called from the HTTP Error Handler.
func setHTTPRespHeaders(_ context.Context, rw http.ResponseWriter, _ protoreflect.ProtoMessage) error {
	for _, headerToDelete := range httpRespHeadersToDelete {
		rw.Header().Del(headerToDelete)
	}
	for headerKeyToAdd, headerValueToAdd := range httpRespHeadersToAdd {
		rw.Header().Set(headerKeyToAdd, headerValueToAdd)
	}
	return nil
}

var corsHeadersToAdd = map[string]string{
	"Access-Control-Allow-Origin":  "*",
	"Access-Control-Allow-Methods": "POST, GET, OPTIONS, PUT, DELETE",
	"Access-Control-Allow-Headers": "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization",
}

var httpRespHeadersToAdd = map[string]string{
	"Content-Security-Policy":   "default-src 'self'",
	"Content-Type":              "application/json",
	"Strict-Transport-Security": "max-age=31536000; includeSubDomains; preload",
	"X-Content-Type-Options":    "nosniff",
	"X-Frame-Options":           "SAMEORIGIN",
	"X-XSS-Protection":          "1; mode=block",
}

var httpRespHeadersToDelete = []string{"Grpc-Metadata-Content-Type"}
