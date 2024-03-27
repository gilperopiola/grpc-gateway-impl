package http

import (
	"context"
	"net/http"

	"github.com/gilperopiola/grpc-gateway-impl/pkg/v1/errs"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"go.uber.org/zap"
	"google.golang.org/grpc/status"
)

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
