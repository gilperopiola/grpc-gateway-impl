package middleware

import (
	"context"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc/status"
)

// httpErrorResponseBody is the struct that gets marshaled onto the HTTP Response body when an error happens.
type httpErrorResponseBody struct {
	Error string `json:"error"`
}

// newHTTPErrorResponse is a helper function that returns a new httpErrorResponseBody.
func newHTTPErrorResponse(message string) httpErrorResponseBody {
	return httpErrorResponseBody{Error: message}
}

// handleHTTPError is a custom error handler for the HTTP Gateway. It's pretty simple.
// It just converts the gRPC error to an HTTP error and writes it to the response.
// There are some special cases based on the HTTP Status.
func handleHTTPError(ctx context.Context, mux *runtime.ServeMux, mar runtime.Marshaler, w http.ResponseWriter, r *http.Request, err error) {

	// Set the Content-Type header.
	w.Header().Set("Content-Type", "application/json")

	var (
		httpRespBuffer          = []byte{}
		grpcStatus              = status.Convert(err)
		httpStatus, httpRespMsg = mapGRPCStatusToHTTPResponseData(grpcStatus)
		httpRespBody            = newHTTPErrorResponse(httpRespMsg)
	)

	// Marshal the error response into a buffer. If it fails, we just write a generic error message
	// and return a 500 Internal Server Error. This is very unlikely to happen, I guess sometimes
	// it's better to be safe than sorry.
	if httpRespBuffer, err = mar.Marshal(httpRespBody); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(marshalErrRespBody))
		return
	}

	// Handle some special cases based on the HTTP Status. If the status is not one of the special cases,
	// the status and the buffer are returned as is.
	httpStatus, httpRespBuffer = handle4XXError(httpStatus, httpRespBuffer, w)

	// Return 5XX errors as 500 Internal Server Error.
	httpStatus, httpRespBuffer = handle5XXError(httpStatus, httpRespBuffer)

	// Write the response.
	w.WriteHeader(httpStatus)
	w.Write(httpRespBuffer)
}

// handle4XXError is a helper function that handles some special error cases based on the HTTP Status.
// If the status is not one of the special cases, it just returns the status and the buffer as is.
func handle4XXError(httpStatus int, buffer []byte, w http.ResponseWriter) (int, []byte) {

	// 401 Unauthorized 						-> We just set a generic message + WWW-Authenticate header.
	// 403 Forbidden 							-> Generic message.
	// 404 Not Found / 405 Method Not Allowed 	-> Generic message + always return 404 Not Found.

	switch httpStatus {
	case http.StatusUnauthorized: // 401
		w.Header().Set("WWW-Authenticate", "Bearer")
		return http.StatusUnauthorized, []byte(unauthorizedErrRespBody)

	case http.StatusForbidden: // 403
		return http.StatusForbidden, []byte(forbiddenErrRespBody)

	case http.StatusNotFound, http.StatusMethodNotAllowed: // 404 / 405
		return http.StatusNotFound, []byte(notFoundErrRespBody)
	}

	// If not, return as is.
	return httpStatus, buffer
}

// handle5XXError is a helper function that handles the 503 error case.
// It just returns 503 as 500 Internal Server Error.
func handle5XXError(httpStatus int, buffer []byte) (int, []byte) {

	// 500 -> We just set a generic message.
	// 503 -> Generic message + always return 500 Internal Server Error.

	switch httpStatus {
	case http.StatusInternalServerError: // 500
		return http.StatusInternalServerError, []byte(internalErrRespBody)

	case http.StatusServiceUnavailable: // 503
		return http.StatusInternalServerError, []byte(svcUnavailErrRespBody)
	}

	// If not, return as is.
	return httpStatus, buffer
}

// mapGRPCStatusToHTTPResponseData returns the HTTP status code and the HTTP error response body
// based on the gRPC status code.
func mapGRPCStatusToHTTPResponseData(grpcStatus *status.Status) (int, string) {
	httpStatus := runtime.HTTPStatusFromCode(grpcStatus.Code())
	return httpStatus, grpcStatus.Message()
}

// These strings are the JSON representations of an httpErrorResponse. It's what gets sent as the response's body when an error occurs.
const (
	marshalErrRespBody      = `{"error": "unexpected error on our end, marshal of the error response failed."}`
	notFoundErrRespBody     = `{"error": "not found, check the docs for the correct path and method."}`
	unauthorizedErrRespBody = `{"error": "unauthorized, please authenticate first."}`
	forbiddenErrRespBody    = `{"error": "forbidden, you aren't allowed to access this resource."}`
	svcUnavailErrRespBody   = `{"error": "the service is unavailable right now, maybe try again later."}`
	internalErrRespBody     = `{"error": "internal server error, something went wrong on our end."}`
)
