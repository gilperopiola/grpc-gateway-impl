package http

import (
	"net/http"
	"time"

	"github.com/gilperopiola/grpc-gateway-impl/app/core"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"go.uber.org/zap"
)

/* ----------------------------------- */
/*         - HTTP Middleware -         */
/* ----------------------------------- */

// Some middleware are passed as ServeMuxOptions when the mux is created.
// Some are wrapped around the Mux afterwards.

// MuxOptionsMiddleware returns all the HTTP middleware that are used as ServeMuxOptions.
func MuxOptionsMiddleware() []runtime.ServeMuxOption {
	return []runtime.ServeMuxOption{
		runtime.WithErrorHandler(HandleHTTPError),
	}
}

// MuxWrapperMiddleware returns the middleware to be wrapped around the HTTP Gateway's Mux.
func MuxWrapperMiddleware() func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return corsMiddleware(
			loggerMiddleware(
				modifyHeadersMiddleware(next),
			),
		)
	}
}

func loggerMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		duration := time.Since(start)

		zap.S().Infow("HTTP Request", core.ZapEndpoint(r.Method+" "+r.URL.Path), core.ZapDuration(duration))

		// Most HTTP logs come with a gRPC log before, as HTTP acts as a gateway to gRPC.
		// As such, we add a new line to separate the logs and easily identify different requests.
		// The only exception would be if there was an error before calling the gRPC handlers.
		zap.L().Info("\n") // T0D0 . Is this ok?
	})
}

// modifyHeadersMiddleware executes before the response is written to the client.
// It's also called from the HTTP Error Handler.
func modifyHeadersMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		next.ServeHTTP(w, r)

		for _, value := range responseHeadersToDelete {
			w.Header().Del(value)
		}
		for key, value := range responseHeadersToAdd {
			w.Header().Set(key, value)
		}
	})
}

var responseHeadersToAdd = map[string]string{
	"Content-Security-Policy":   "default-src 'self'",
	"Content-Type":              "application/json",
	"Strict-Transport-Security": "max-age=31536000; includeSubDomains; preload",
	"X-Content-Type-Options":    "nosniff",
	"X-Frame-Options":           "SAMEORIGIN",
	"X-XSS-Protection":          "1; mode=block",
}

var responseHeadersToDelete = []string{"Grpc-Metadata-Content-Type"}

// corsMiddleware adds CORS headers to the response and handles preflight requests.
func corsMiddleware(next http.Handler) http.Handler {
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

var corsHeadersToAdd = map[string]string{
	"Access-Control-Allow-Origin":  "*",
	"Access-Control-Allow-Methods": "POST, GET, OPTIONS, PUT, DELETE",
	"Access-Control-Allow-Headers": "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization",
}
