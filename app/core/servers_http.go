package core

import (
	"context"
	"net/http"
	"time"

	"github.com/gilperopiola/grpc-gateway-impl/app/core/errs"
	"github.com/gilperopiola/grpc-gateway-impl/app/core/pbs"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
)

/* -~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~- */
/*          - HTTP Gateway -           */
/* -~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~- */

func newHTTPGateway(serveOpts []runtime.ServeMuxOption, middleware func(http.Handler) http.Handler, grpcDialOpts []grpc.DialOption) *http.Server {
	mux := runtime.NewServeMux(serveOpts...)

	if err := pbs.RegisterUsersServiceHandlerFromEndpoint(context.Background(), mux, GRPCPort, grpcDialOpts); err != nil {
		LogUnexpectedAndPanic(err)
	}

	return &http.Server{
		Addr:    HTTPPort,
		Handler: middleware(mux),
	}
}

func runHTTP(httpGateway *http.Server) {
	zap.S().Infof("GRPC Gateway Implementation | HTTP Port %s 🚀", HTTPPort)

	go func() {
		if err := httpGateway.ListenAndServe(); err != http.ErrServerClosed {
			LogUnexpectedAndPanic(err)
		}
	}()
}

func shutdownHTTP(httpGateway *http.Server) {
	zap.S().Info("Shutting down HTTP server...")

	timeout := 4 * time.Second
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	if err := httpGateway.Shutdown(ctx); err != nil {
		LogUnexpectedAndPanic(err)
	}
}

/* -~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~- */
/*         - HTTP Middleware -         */
/* -~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~- */

// Some middleware are passed as ServeMuxOptions when the mux is created.
// Some are wrapped around the Mux afterwards.

// Returns the middleware to be wrapped around the HTTP Gateway's Mux.
func defaultHTTPMiddleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return corsMiddleware(
			requestsLoggerMiddleware(
				setResponseHeadersMiddleware(next),
			),
		)
	}
}

// Returns all the HTTP middleware that are used as ServeMuxOptions.
func defaultHTTPServeOpts() []runtime.ServeMuxOption {
	return []runtime.ServeMuxOption{
		runtime.WithErrorHandler(handleHTTPError),
	}
}

// Logs HTTP Requests.
func requestsLoggerMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		duration := time.Since(start)

		zap.S().Infow("HTTP Request", ZapRoute(r.Method+" "+r.URL.Path), ZapDuration(duration))

		// Most HTTP logs come with a gRPC log before, as HTTP acts as a gateway to gRPC.
		// As such, we add a new line to separate the logs and easily identify different requests.
		// The only exception would be if there was an error before calling the gRPC handlers.
		zap.L().Info("\n") // T0D0 . Is this ok?
	})
}

// Executes before response headers are fully set.
// Also called from our HTTP Error Handler Func.
func setResponseHeadersMiddleware(next http.Handler) http.Handler {
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
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		next.ServeHTTP(w, r)

		for _, value := range headersToDelete {
			w.Header().Del(value)
		}
		for key, value := range headersToAdd {
			w.Header().Set(key, value)
		}
	})
}

// Adds CORS headers to the response and handles preflight requests.
func corsMiddleware(next http.Handler) http.Handler {
	var headersToAdd = map[string]string{
		"Access-Control-Allow-Origin":  "*",
		"Access-Control-Allow-Methods": "POST, GET, OPTIONS, PUT, DELETE",
		"Access-Control-Allow-Headers": "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization",
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		for key, value := range headersToAdd {
			w.Header().Set(key, value)
		}
		if r.Method == "OPTIONS" { // Preflight request
			w.WriteHeader(http.StatusOK)
			return
		}
		next.ServeHTTP(w, r)
	})
}

/* -~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~- */
/*        - HTTP Error Handler -       */
/* -~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~- */

// handleHTTPError is a custom error handler for the HTTP Gateway. It's pretty simple.
// It converts the gRPC error to an HTTP error and writes it to the response.
func handleHTTPError(ctx context.Context, mux *runtime.ServeMux, m runtime.Marshaler, w http.ResponseWriter, _ *http.Request, err error) {
	grpcStatus := status.Convert(err)
	httpStatus := runtime.HTTPStatusFromCode(grpcStatus.Code())

	// This function stops the execution chain, so we manually call the forwardResponseOptions to set the headers.
	for _, forwardResponseFn := range mux.GetForwardResponseOptions() {
		forwardResponseFn(ctx, w, nil)
	}

	// Create and marshal an httpError into a []byte buffer. If it fails (unlikely), we return 500 Internal Server Error.
	var httpBody []byte
	if httpBody, err = m.Marshal(httpError{grpcStatus.Message()}); err != nil {
		LogUnexpected(err)
		httpStatus = http.StatusInternalServerError
	}

	httpStatus, httpBody = standardizeErrorResponse(w, httpStatus, httpBody)
	w.WriteHeader(httpStatus)
	w.Write(httpBody)
}

// httpError is the struct that gets marshalled onto the HTTP Response body when an error happens.
// This is what the client would see. The format is '{"error": "error message."}'.
// If this format is changed, then the ErrBodies in errs.go should also change.
type httpError struct {
	Error string `json:"error"`
}

// standardizeErrorResponse returns a generic HTTP Status Code and Body for each HTTP Status Code.
func standardizeErrorResponse(w http.ResponseWriter, status int, respBody []byte) (int, []byte) {
	switch status {
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
	return status, respBody
}
