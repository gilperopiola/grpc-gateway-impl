package servers

import (
	"fmt"
	"net/http"

	"github.com/gilperopiola/god"
	"github.com/gilperopiola/grpc-gateway-impl/app/core"
	"github.com/gilperopiola/grpc-gateway-impl/app/core/errs"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc/status"
)

func handleHTTPError(_ god.Ctx, _ *runtime.ServeMux, _ runtime.Marshaler, rw http.ResponseWriter, _ *http.Request, err error) {
	grpcStatus := status.Convert(err)
	httpStatus := runtime.HTTPStatusFromCode(grpcStatus.Code())
	httpRespBody := grpcStatus.Message()

	finalizeErrorResponse(rw.Header().Set, httpStatus, &httpRespBody)

	deleteGRPCHeader(rw)
	rw.WriteHeader(httpStatus)
	rw.Write([]byte(httpRespBody))
}

// Called just before sending back an HTTP error response.
// Touches up the outgoing response body and headers a bit based on the status code.
// It's like a middleware for our HTTP error handling middleware.
func finalizeErrorResponse(setHeaderFn func(string, string), status int, body *string) {
	tempBody := *body

	switch status {

	// HTTP 400.
	// Return as is.
	case http.StatusBadRequest:
		break

	// HTTP 401.
	// Add Header + Generic response.
	case http.StatusUnauthorized:
		setHeaderFn("WWW-Authenticate", "Bearer")
		tempBody = errs.HTTPUnauthorized

	// HTTP 403.
	// Generic response.
	case http.StatusForbidden:
		tempBody = errs.HTTPForbidden

	// HTTP 404 ~ 405.
	// Generic response if error is GRPC Gateway's Not-Found.
	// Return as is otherwise.
	//
	// This is because our Service might return Not Found
	// but with a more specific error message.
	case http.StatusNotFound, http.StatusMethodNotAllowed:
		if *body == `{"error": "Not Found"}` {
			tempBody = errs.HTTPRouteNotFound
		}

	// HTTP 409.
	// Return as is.
	case http.StatusConflict:
		break

	// HTTP 500.
	// Log + Generic response.
	case http.StatusInternalServerError:
		core.LogStrange("HTTP Error 500: " + tempBody)
		tempBody = errs.HTTPInternal

	// HTTP 503.
	// Log + Return as is.
	//
	// Failed health checks return 503.
	case http.StatusServiceUnavailable:
		core.LogStrange("HTTP Error 503: " + tempBody)

	default:
		core.LogStrange(fmt.Sprintf("HTTP Error %d (unhandled): %s", status, tempBody))
	}

	// For all cases, we set this as the response body format.
	*body = fmt.Sprintf(`{"error": "%s"}`, tempBody)
}
