package http

import (
	"context"
	"net/http"

	"github.com/gilperopiola/grpc-gateway-impl/pkg/v1/dependencies"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/protobuf/reflect/protoreflect"
)

/* ----------------------------------- */
/*         - HTTP Middleware -         */
/* ----------------------------------- */
// Some middleware are passed as ServeMuxOptions when the mux is created.
// Some are wrapped around the mux afterwards.

// MuxWrapperFunc is a middleware that wraps around the HTTP Gateway's mux.
type MuxWrapperFunc func(next http.Handler) http.Handler

// AllMiddleware returns all the HTTP middleware that are used as ServeMuxOptions.
func AllMiddleware() []runtime.ServeMuxOption {
	return []runtime.ServeMuxOption{
		runtime.WithErrorHandler(handleHTTPErr), // Stops other middleware if an error happens.
		runtime.WithForwardResponseOption(setHTTPResponseHeaders),
	}
}

// AllMiddlewareWrapper returns the middleware to be wrapped around the HTTP Gateway's Mux.
func AllMiddlewareWrapper(logger *dependencies.Logger) MuxWrapperFunc {
	sugar := logger.Sugar()
	return func(next http.Handler) http.Handler {
		return handleCORS(
			logger.LogHTTP(next, sugar),
		)
	}
}

// handleCORS adds CORS headers to the response and handles preflight requests.
func handleCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		// Add headers.
		for key, value := range corsHeaders {
			w.Header().Set(key, value)
		}

		// Handle preflight requests.
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		// Pass down the chain to next handler if not OPTIONS.
		next.ServeHTTP(w, r)
	})
}

var corsHeaders = map[string]string{
	"Access-Control-Allow-Origin":  "*",
	"Access-Control-Allow-Methods": "POST, GET, OPTIONS, PUT, DELETE",
	"Access-Control-Allow-Headers": "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization",
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

// setHTTPResponseHeadersWrapper allows setHTTPResponseHeaders to be called without context and message.
func setHTTPResponseHeadersWrapper(rw http.ResponseWriter) error {
	return setHTTPResponseHeaders(context.Background(), rw, protoreflect.ProtoMessage(nil))
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
