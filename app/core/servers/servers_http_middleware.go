package servers

import (
	"net/http"

	"github.com/gilperopiola/grpc-gateway-impl/app/core"
)

/* -~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~- */
/*         - HTTP Middleware -         */
/* -~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~- */

type middlewareFunc func(http.Handler) http.Handler

// Returns the middleware to be wrapped around the HTTP Gateway's Mux.
func defaultHTTPMiddleware() func(http.Handler) http.Handler {
	return func(handler http.Handler) http.Handler {
		return corsMiddleware(
			core.LogHTTPRequest(
				setResponseHeadersMiddleware(handler),
			),
		)
	}
}

// Executes before response headers are fully set.
// Also called from our HTTP Error Handler Func.
func setResponseHeadersMiddleware(handler http.Handler) http.Handler {
	var (
		headersToDelete = []string{"Grpc-Metadata-Content-Type"}
		headersToAdd    = map[string]string{
			"Content-Security-Policy":   "default-src 'self'",
			"Content-Type":              "application/json",
			"Strict-Transport-Security": "max-age=31536000; includeSubDomains; preload",
			"X-Content-Type-Options":    "nosniff",
			"X-Frame-Options":           "SAMEORIGIN",
			"X-XSS-Protection":          "1; mode=block",
		}
	)
	return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		handler.ServeHTTP(rw, req)

		for _, value := range headersToDelete {
			rw.Header().Del(value)
		}
		for key, value := range headersToAdd {
			rw.Header().Set(key, value)
		}
	})
}

// Adds CORS headers to the response and handles preflight requests.
func corsMiddleware(handler http.Handler) http.Handler {
	headersToAdd := map[string]string{
		"Access-Control-Allow-Origin":  "*",
		"Access-Control-Allow-Methods": "POST, GET, OPTIONS, PUT, DELETE",
		"Access-Control-Allow-Headers": "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization",
	}

	return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		for key, value := range headersToAdd {
			rw.Header().Set(key, value)
		}
		if req.Method == "OPTIONS" { // Preflight
			rw.WriteHeader(http.StatusOK)
			return
		}
		handler.ServeHTTP(rw, req)
	})
}
