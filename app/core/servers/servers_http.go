package servers

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gilperopiola/grpc-gateway-impl/app/core"
	"github.com/gilperopiola/grpc-gateway-impl/app/core/errs"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc/status"
)

/* -~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~- */
/*          - HTTP Gateway -           */
/* -~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~- */

type middlewareFunc func(http.Handler) http.Handler

// Returns the middleware to be wrapped around the HTTP Gateway's Mux.
func getHTTPMiddleware() middlewareFunc {
	return func(handler http.Handler) http.Handler {
		return handleCORS(core.LogHTTPRequest(setResponseHeaders(handler)))
	}
}

// Middleware. Adds CORS headers to the response and handles preflight requests
var handleCORS middlewareFunc = func(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		core.LogIfDebug("Handling CORS | " + req.Method + " request from " + req.RemoteAddr)

		for key, value := range corsHeaders {
			rw.Header().Set(key, value)
		}
		if req.Method == "OPTIONS" { // Preflight
			rw.WriteHeader(http.StatusOK)
			return
		}
		handler.ServeHTTP(rw, req)
	})
}

// Middleware
var setResponseHeaders middlewareFunc = func(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		handler.ServeHTTP(rw, req)
		for _, value := range deleteHeaders {
			rw.Header().Del(value)
		}
		for key, value := range defaultHeaders {
			rw.Header().Set(key, value)
		}
	})
}

/* -~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~- */

// ServeMuxOptions are applied to the HTTP Gateway's Mux on creation.
// For now there's only an error handler.

// Returns our ServeMuxOptions.
func getHTTPServeMuxOptions() []runtime.ServeMuxOption {
	return []runtime.ServeMuxOption{
		runtime.WithErrorHandler(handleHTTPError),
		runtime.WithRoutingErrorHandler(func(c core.Ctx, m *runtime.ServeMux, marsh runtime.Marshaler, rw http.ResponseWriter, req *http.Request, idk int) {
			core.LogWeirdBehaviour("Route not found")
		}),
	}
}

func handleHTTPError(c core.Ctx, mux *runtime.ServeMux, m runtime.Marshaler, rw http.ResponseWriter, _ *http.Request, err error) {

	// HTTP Errors stop the HTTP Middleware execution chain, so we call forwardResponseOptions to set the headers.
	for _, forwardResponseFn := range mux.GetForwardResponseOptions() {
		forwardResponseFn(c, rw, nil)
	}

	grpcStatus := status.Convert(err)
	httpStatus := runtime.HTTPStatusFromCode(grpcStatus.Code())
	httpBody := newHTTPErrorRespBody(grpcStatus.Message())

	updateHTTPErrorResponse(rw, httpStatus, &httpBody)

	rw.WriteHeader(httpStatus)
	rw.Write([]byte(httpBody))
}

// Modifies the HTTP Error response body and headers based on the HTTP Status.
func updateHTTPErrorResponse(rw http.ResponseWriter, status int, body *string) {
	switch status {
	case http.StatusBadRequest:
		// Return as is
	case http.StatusUnauthorized:
		rw.Header().Set("WWW-Authenticate", "Bearer")
		*body = newHTTPErrorRespBody(errs.HTTPUnauthorized)
	case http.StatusForbidden:
		*body = newHTTPErrorRespBody(errs.HTTPForbidden)
	case http.StatusNotFound, http.StatusMethodNotAllowed:
		*body = newHTTPErrorRespBody(errs.HTTPNotFound)
	case http.StatusConflict:
		*body = newHTTPErrorRespBody(errs.HTTPConflict)
	case http.StatusInternalServerError:
		*body = newHTTPErrorRespBody(errs.HTTPInternal)
	case http.StatusServiceUnavailable:
		*body = newHTTPErrorRespBody(errs.HTTPUnavailable)
	default:
		core.LogWeirdBehaviour("HTTP Error Status: " + strconv.Itoa(status))
	}
}

func newHTTPErrorRespBody(msg string) string {
	return fmt.Sprintf(`{"error": "%s"}`, msg)
}

/* -~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~- */

var (
	corsHeaders = map[string]string{
		"Access-Control-Allow-Origin":  "*",
		"Access-Control-Allow-Methods": "POST, GET, OPTIONS, PUT, DELETE",
		"Access-Control-Allow-Headers": "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization",
	}
	defaultHeaders = map[string]string{
		"Content-Type":              "application/json",
		"Content-Security-Policy":   "default-src 'self'",
		"Strict-Transport-Security": "max-age=31536000; includeSubDomains; preload",
		"X-Content-Type-Options":    "nosniff",
		"X-Frame-Options":           "SAMEORIGIN",
		"X-XSS-Protection":          "1; mode=block",
	}
	deleteHeaders = []string{"Grpc-Metadata-Content-Type"}
)
