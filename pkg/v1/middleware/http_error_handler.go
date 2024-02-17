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

// handleHTTPError is a custom error handler for the HTTP Gateway. It's pretty simple.
// It just converts the gRPC error to an HTTP error and writes it to the response.
// There are some special cases based on the HTTP Status.
func handleHTTPError(ctx context.Context, mux *runtime.ServeMux, mar runtime.Marshaler, w http.ResponseWriter, r *http.Request, err error) {

	// Set the Content-Type header.
	w.Header().Set("Content-Type", "application/json")

	var (
		buffer                   = []byte{}
		grpcStatus               = status.Convert(err)
		httpStatus, responseBody = getHTTPResponseStatusAndBody(grpcStatus)
	)

	// Marshal the error response into a buffer. If it fails, we just write a generic error message
	// and return a 500 Internal Server Error.
	if buffer, err = mar.Marshal(responseBody); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(marshalErrorResponseBody))
		return
	}

	// Handle some special cases based on the HTTP Status. If the status is not one of the special cases,
	// the status and the buffer are returned as is.
	httpStatus, buffer = handle4XXError(httpStatus, buffer, w)

	// Write the response.
	w.WriteHeader(httpStatus)
	w.Write(buffer)
}

// handle4XXError is a helper function that handles some special error cases based on the HTTP Status.
// If the status is not one of the special cases, it just returns the status and the buffer as is.
func handle4XXError(httpStatus int, buffer []byte, w http.ResponseWriter) (int, []byte) {

	// 401 (we just set a generic message + WWW-Authenticate header).
	if httpStatus == http.StatusUnauthorized {
		w.Header().Set("WWW-Authenticate", "Bearer")
		return http.StatusUnauthorized, []byte(unauthorizedErrorResponseBody)
	}

	// 403 (generic message).
	if httpStatus == http.StatusForbidden {
		return http.StatusForbidden, []byte(forbiddenErrorResponseBody)
	}

	// 404 / 405 (generic message + always return 404).
	if httpStatus == http.StatusNotFound || httpStatus == http.StatusMethodNotAllowed {
		return http.StatusNotFound, []byte(notFoundErrorResponseBody)
	}

	return httpStatus, buffer
}

// getHTTPResponseStatusAndBody returns the HTTP status and the HTTP error response body.
func getHTTPResponseStatusAndBody(grpcStatus *status.Status) (int, httpErrorResponseBody) {
	httpStatus := runtime.HTTPStatusFromCode(grpcStatus.Code())
	httpResponse := httpErrorResponseBody{Error: grpcStatus.Message()}
	return httpStatus, httpResponse
}

// These strings are the JSON representations of an httpErrorResponse. It's what gets sent as the request body when an error occurs.
const (
	marshalErrorResponseBody      = `{"error": "failed to marshal error response"}`
	notFoundErrorResponseBody     = `{"error": "not found, check the docs for the correct path and method"}`
	unauthorizedErrorResponseBody = `{"error": "unauthorized"}`
	forbiddenErrorResponseBody    = `{"error": "forbidden"}`
)
