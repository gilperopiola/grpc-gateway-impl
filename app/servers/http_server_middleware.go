package servers

import (
	"context"
	"net/http"
	"time"

	"github.com/gilperopiola/grpc-gateway-impl/app/core"
	"github.com/gilperopiola/grpc-gateway-impl/app/core/errs"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"go.uber.org/zap"
	"google.golang.org/grpc/status"
)

/* ----------------------------------- */
/*         - HTTP Middleware -         */
/* ----------------------------------- */

// Some middleware are passed as ServeMuxOptions when the mux is created.
// Some are wrapped around the Mux afterwards.

// ServeMuxOpts returns all the HTTP middleware that are used as ServeMuxOptions.
func ServeMuxOpts() []runtime.ServeMuxOption {
	return []runtime.ServeMuxOption{
		runtime.WithErrorHandler(HandleHTTPError),
	}
}

// MiddlewareWrapper returns the middleware to be wrapped around the HTTP Gateway's Mux.
func MiddlewareWrapper() func(next http.Handler) http.Handler {
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
		next.ServeHTTP(w, r) // Call next handler.
	})
}

var corsHeadersToAdd = map[string]string{
	"Access-Control-Allow-Origin":  "*",
	"Access-Control-Allow-Methods": "POST, GET, OPTIONS, PUT, DELETE",
	"Access-Control-Allow-Headers": "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization",
}

/* ----------------------------------- */
/*        - HTTP Error Handler -       */
/* ----------------------------------- */

// HandleHTTPError is a custom error handler for the HTTP Gateway. It's pretty simple.
// It converts the gRPC error to an HTTP error and writes it to the response.
func HandleHTTPError(ctx context.Context, mux *runtime.ServeMux, m runtime.Marshaler, w http.ResponseWriter, _ *http.Request, err error) {
	grpcStatus := status.Convert(err)                           // err 			-> gRPC Status.
	httpStatus := runtime.HTTPStatusFromCode(grpcStatus.Code()) // gRPC Status -> HTTP Status.

	// This function stops the execution chain, so we manually call the forwardResponseOptions to set the headers.
	for _, fn := range mux.GetForwardResponseOptions() {
		fn(ctx, w, nil)
	}

	// Create and marshal an httpError into a []byte buffer. If it fails (unlikely), we return 500 Internal Server Error.
	var httpBody []byte
	if httpBody, err = m.Marshal(httpError{Error: grpcStatus.Message()}); err != nil {
		httpStatus = http.StatusInternalServerError
		zap.S().Errorf("Failed to marshal error message: %v", err)
	}

	// We send generic responses for some HTTP Status Codes.
	httpStatus, httpBody = getGenericErrorResponse(w, httpStatus, httpBody)

	w.WriteHeader(httpStatus)
	w.Write(httpBody)
}

// getGenericErrorResponse returns a generic HTTP Status Code and Body for each HTTP Status Code.
func getGenericErrorResponse(w http.ResponseWriter, httpStatus int, httpBody []byte) (int, []byte) {
	switch httpStatus {
	case http.StatusUnauthorized:
		w.Header().Set("WWW-Authenticate", "Bearer")
		return http.StatusUnauthorized, []byte(errs.HTTPUnauthorizedErrBody) // ------------ 401 (WWW-Authenticate: Bearer)
	case http.StatusForbidden:
		return http.StatusForbidden, []byte(errs.HTTPForbiddenErrBody) // ------------------ 403 (Forbidden)
	case http.StatusNotFound, http.StatusMethodNotAllowed:
		return http.StatusNotFound, []byte(errs.HTTPNotFoundErrBody) // -------------------- 404 & 405 (Not Found)
	case http.StatusInternalServerError:
		return http.StatusInternalServerError, []byte(errs.HTTPInternalErrBody) // --------- 500 (Internal Server Error)
	case http.StatusServiceUnavailable:
		return http.StatusInternalServerError, []byte(errs.HTTPServiceUnavailErrBody) // --- 503 (Service Unavailable)
	}
	return httpStatus, httpBody
}

// httpError is the struct that gets marshalled onto the HTTP Response body when an error happens.
// This is what the client would see. The format is '{"error": "error message."}'.
// If this format is changed, then the ErrBodies in errs.go should also change.
type httpError struct {
	Error string `json:"error"`
}
