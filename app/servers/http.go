package servers

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gilperopiola/god"
	"github.com/gilperopiola/grpc-gateway-impl/app/core"
	"github.com/gilperopiola/grpc-gateway-impl/app/core/errs"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/reflect/protoreflect"
)

/* -~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~- */
/*          - HTTP Gateway -           */
/* -~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~- */

// -> Middleware

type middlewareFunc func(http.Handler) http.Handler

// Returns the middleware to be wrapped around the HTTP Gateway's Mux
func getHTTPMiddlewareChain() middlewareFunc {
	return func(handler http.Handler) http.Handler {
		return addCustomRespWriter(
			handleCORS(
				setResponseHeaders(
					core.LogHTTPRequest(handler),
				),
			),
		)
	}
}

// Replaces the default ResponseWriter with our CustomResponseWriter
var addCustomRespWriter middlewareFunc = func(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		handler.ServeHTTP(&core.CustomResponseWriter{rw, http.StatusOK}, req)
	})
}

// Adds CORS headers and handles preflight requests
var handleCORS middlewareFunc = func(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		core.LogIfDebug("CORS " + req.Method + " from " + req.RemoteAddr)

		for key, value := range corsHeaders {
			rw.Header().Set(key, value)
		}

		if req.Method == "OPTIONS" {
			// When the request is a preflight, the client is asking for permission to make the actual request.
			// We respond with the allowed methods and headers.
			rw.WriteHeader(http.StatusOK)
			return
		}

		handler.ServeHTTP(rw, req)
	})
}

// Sets default headers and removes the GRPC header
var setResponseHeaders middlewareFunc = func(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		for key, value := range defaultHeaders {
			rw.Header().Set(key, value)
		}
		handler.ServeHTTP(rw, req)
		deleteGRPCHeader(rw)
	})
}

/* -~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~- */

// -> Serve Mux Options

// Returns our ServeMuxOptions.
// ServeMuxOptions are applied to the HTTP Gateway's Mux on creation.
// For now there's only an error handler.
func getHTTPMuxOpts() []runtime.ServeMuxOption {
	return []runtime.ServeMuxOption{
		runtime.WithErrorHandler(handleHTTPError),
		runtime.WithForwardResponseOption(func(_ god.Ctx, rw http.ResponseWriter, _ protoreflect.ProtoMessage) error {
			deleteGRPCHeader(rw)
			return nil
		}),
	}
}

func handleHTTPError(_ god.Ctx, _ *runtime.ServeMux, _ runtime.Marshaler, rw http.ResponseWriter, _ *http.Request, err error) {
	grpcStatus := status.Convert(err)
	httpStatus := runtime.HTTPStatusFromCode(grpcStatus.Code())
	httpRespBody := grpcStatus.Message()

	finalizeErrorResponse(rw, httpStatus, &httpRespBody)

	deleteGRPCHeader(rw)
	rw.WriteHeader(httpStatus)
	rw.Write([]byte(httpRespBody))
}

func deleteGRPCHeader(rw http.ResponseWriter) {
	rw.Header().Del(grpcHeader)
}

// Modifies the HTTP Error response body and headers based on the HTTP Status.
func finalizeErrorResponse(rw http.ResponseWriter, status int, body *string) {
	errMsg := *body

	switch status {
	case http.StatusBadRequest:
		// Return as is.

	case http.StatusUnauthorized:
		rw.Header().Set("WWW-Authenticate", "Bearer")
		errMsg = errs.HTTPUnauthorized

	case http.StatusForbidden:
		errMsg = errs.HTTPForbidden

	case http.StatusNotFound, http.StatusMethodNotAllowed:
		// GRPC Gateway returns this ⬇️ for not-found routes.
		if errMsg == `{"error": "Not Found"}` {
			errMsg = errs.HTTPRouteNotFound
		}
		// If we get a not-found from the service, we return the message as is.

	case http.StatusConflict:
		errMsg = errs.HTTPConflict

	case http.StatusInternalServerError:
		core.LogWeirdBehaviour("HTTP Error 500: " + errMsg)
		errMsg = errs.HTTPInternal

	case http.StatusServiceUnavailable:
		core.LogWeirdBehaviour("HTTP Error 503: " + errMsg)
		errMsg = errs.HTTPUnavailable

	default:
		core.LogWeirdBehaviour("HTTP Error Status: " + strconv.Itoa(status))
	}

	*body = fmt.Sprintf(`{"error": "%s"}`, errMsg)
}

/* -~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~- */

var (
	defaultHeaders = map[string]string{
		"Content-Type":              "application/json",
		"Content-Security-Policy":   "default-src 'self'",
		"Strict-Transport-Security": "max-age=31536000; includeSubDomains; preload",
		"X-Content-Type-Options":    "nosniff",
		"X-Frame-Options":           "SAMEORIGIN",
		"X-XSS-Protection":          "1; mode=block",
	}

	corsHeaders = map[string]string{
		"Access-Control-Allow-Origin":  "*",
		"Access-Control-Allow-Methods": "POST, GET, OPTIONS, PUT, DELETE",
		"Access-Control-Allow-Headers": "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization",
	}

	// Deleted in HTTP Response
	grpcHeader = "Grpc-Metadata-Content-Type"
)
