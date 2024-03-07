package http

import (
	"context"
	"net/http"

	"github.com/gilperopiola/grpc-gateway-impl/pkg/v1/errs"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/reflect/protoreflect"
)

/* ----------------------------------- */
/*        - HTTP Error Handler -       */
/* ----------------------------------- */

// httpErrorResponseBody is the struct that gets marshalled onto the HTTP Response body when an error happens.
// The format is basically '{"error": "error message."}'.
type httpErrorResponseBody struct {
	Error string `json:"error"`
}

// handleHTTPError is a custom error handler for the HTTP Gateway. It's pretty simple.
// It just converts the gRPC error to an HTTP error and writes it to the response.
// There are some special cases based on the HTTP Status.
func handleHTTPError(ctx context.Context, mux *runtime.ServeMux, mar runtime.Marshaler, w http.ResponseWriter, r *http.Request, err error) {

	// First, get gRPC status from the error. Then, derive the HTTP Response info from that status.
	var (
		grpcStatus       = status.Convert(err)
		httpStatus       = runtime.HTTPStatusFromCode(grpcStatus.Code())
		httpResponseBody = httpErrorResponseBody{Error: grpcStatus.Message()}
	)

	// Marshal the error into a buffer. If it fails (unlikely), we just return a generic 500 Internal Server Error.
	var outBuffer []byte
	if outBuffer, err = mar.Marshal(httpResponseBody); err != nil {
		httpStatus = http.StatusInternalServerError
	}

	// If the HTTP Response Code is 4xx or 5xx, we change the response a little bit, mainly sending generic messages.
	// Otherwise, the status and the buffer are returned as is.
	httpStatus, outBuffer = handle4xxOr5xxError(httpStatus, outBuffer, w)

	// Write the response.
	setHTTPResponseHeaders(ctx, w, protoreflect.ProtoMessage(nil))
	w.WriteHeader(httpStatus)
	w.Write(outBuffer)
}

// handle4xxOr5xxError is a helper function that handles some special error cases based on the HTTP Status.
// If the status is not one of the special cases, it just returns the status and the buffer as is.
func handle4xxOr5xxError(httpStatus int, buffer []byte, w http.ResponseWriter) (int, []byte) {

	// 401 Unauthorized 												-> Generic message + WWW-Authenticate header.
	// 403 Forbidden 														-> Generic message.
	// 404 Not Found / 405 Method Not Allowed 	-> Generic message + always return 404 Not Found.
	// 500 Internal Server Error								-> Generic message.
	// 503 Service Unavailable			 						-> Generic message + always return 500 Internal Server Error.

	switch httpStatus {
	case http.StatusUnauthorized:
		w.Header().Set("WWW-Authenticate", "Bearer")
		return http.StatusUnauthorized, []byte(errs.HTTPUnauthorizedErrBody) // 401

	case http.StatusForbidden:
		return http.StatusForbidden, []byte(errs.HTTPForbiddenErrBody) // 403

	case http.StatusNotFound, http.StatusMethodNotAllowed:
		return http.StatusNotFound, []byte(errs.HTTPNotFoundErrBody) // 404 / 405

	case http.StatusInternalServerError:
		return http.StatusInternalServerError, []byte(errs.HTTPInternalErrBody) // 500

	case http.StatusServiceUnavailable:
		return http.StatusInternalServerError, []byte(errs.HTTPServiceUnavailErrBody) // 503
	}

	// If not, return as is.
	return httpStatus, buffer
}
