package servers

import (
	"fmt"
	"net/http"

	"github.com/gilperopiola/god"
	"github.com/gilperopiola/grpc-gateway-impl/app/core"
	"github.com/gilperopiola/grpc-gateway-impl/app/core/errs"
	"github.com/gilperopiola/grpc-gateway-impl/app/tools"

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
		handler.ServeHTTP(tools.NewHTTPRespWriter(rw), req)
	})
}

// Adds CORS headers and handles preflight requests
var handleCORS middlewareFunc = func(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		core.LogDebug("CORS " + req.Method + " from " + req.RemoteAddr)

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

func deleteGRPCHeader(rw http.ResponseWriter) {
	rw.Header().Del(grpcHeader)
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

// Executed just before sending back an HTTP error response.
// Modifies the response body and headers based on the status.
func finalizeErrorResponse(rw http.ResponseWriter, status int, body *string) {
	finalBody := *body

	switch status {

	// 400.
	// Return as is.
	case http.StatusBadRequest:
		break

	// 401.
	// Add Header + Generic response.
	case http.StatusUnauthorized:
		rw.Header().Set("WWW-Authenticate", "Bearer")
		finalBody = errs.HTTPUnauthorized

	// 403.
	// Generic response.
	case http.StatusForbidden:
		finalBody = errs.HTTPForbidden

	// 404 ~ 405.
	// Generic response if error is GRPC Gateway's Not-Found.
	// Return as is otherwise.
	//
	// This is because our Service might return Not Found
	// but with a more specific error message.
	case http.StatusNotFound, http.StatusMethodNotAllowed:
		if *body == `{"error": "Not Found"}` {
			finalBody = errs.HTTPRouteNotFound
		}

	// 409.
	// Return as is.
	case http.StatusConflict:
		break

	// 500.
	// Log + Generic response.
	case http.StatusInternalServerError:
		core.LogWeirdBehaviour("HTTP Error 500: " + finalBody)
		finalBody = errs.HTTPInternal

	// 503.
	// Log + Return as is.
	//
	// Failed health checks return 503.
	case http.StatusServiceUnavailable:
		core.LogWeirdBehaviour("HTTP Error 503: " + finalBody)

	default:
		core.LogWeirdBehaviour(fmt.Sprintf("HTTP Error %d (unhandled): %s", status, finalBody))
	}

	*body = fmt.Sprintf(`{"error": "%s"}`, finalBody)
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
