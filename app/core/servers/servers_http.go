package servers

import (
	"context"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gilperopiola/grpc-gateway-impl/app/core"
	"github.com/gilperopiola/grpc-gateway-impl/app/core/errs"
	"github.com/gilperopiola/grpc-gateway-impl/app/core/pbs"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
)

/* -~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~- */
/*          - HTTP Gateway -           */
/* -~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~- */

func newHTTPGateway(serveMuxOpts []runtime.ServeMuxOption, middlewareFn middlewareFunc, grpcDialOpts []grpc.DialOption) *http.Server {
	mux := runtime.NewServeMux(serveMuxOpts...)
	err := pbs.RegisterUsersServiceHandlerFromEndpoint(context.Background(), mux, core.GRPCPort, grpcDialOpts)
	core.LogPanicIfErr(err)

	return &http.Server{
		Addr:    core.HTTPPort,
		Handler: middlewareFn(mux),
		// TLSConfig: core.GetTLSConfig(core.GetCertPool(core.CertPath), core.GetServerName(core.ServerName)),
	}
}

/* -~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~- */

type middlewareFunc func(http.Handler) http.Handler

// Returns the middleware to be wrapped around the HTTP Gateway's Mux.
func defaultHTTPMiddleware() middlewareFunc {
	return func(handler http.Handler) http.Handler {
		return handleCORS(core.LogHTTPRequest(setResponseHeaders(handler)))
	}
}

// Middleware.
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

// Middleware. Adds CORS headers to the response and handles preflight requests.
var handleCORS middlewareFunc = func(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
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

/* -~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~- */

// ServeMuxOptions are applied to the HTTP Gateway's Mux on creation.
// For now there's only an error handler.

// Returns our ServeMuxOptions.
func defaultHTTPServeMuxOpts() []runtime.ServeMuxOption {
	return []runtime.ServeMuxOption{runtime.WithErrorHandler(handleHTTPError)}
}

func handleHTTPError(c context.Context, mux *runtime.ServeMux, m runtime.Marshaler, rw http.ResponseWriter, _ *http.Request, err error) {

	// HTTP Errors stop the HTTP Middleware execution chain, so we call forwardResponseOptions to set the headers.
	for _, forwardResponseFn := range mux.GetForwardResponseOptions() {
		forwardResponseFn(c, rw, nil)
	}

	grpcStatus := status.Convert(err)
	httpStatus := runtime.HTTPStatusFromCode(grpcStatus.Code())
	httpBody := newHTTPErrorRespBody(grpcStatus.Message())

	modifyHTTPErrorBodyOrHeaders(rw, httpStatus, &httpBody)

	rw.WriteHeader(httpStatus)
	rw.Write([]byte(httpBody))
}

// Modifies the HTTP Error response body and headers based on the HTTP Status.
func modifyHTTPErrorBodyOrHeaders(rw http.ResponseWriter, status int, body *string) {
	switch status {
	case http.StatusBadRequest:
		// Do nothing.
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
