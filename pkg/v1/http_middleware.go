package v1

import (
	"context"
	"io"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/reflect/protoreflect"
)

const (
	contentTypeHeader = "Content-Type"
	contentTypeJSON   = "application/json"
)

/* ----------------------------------- */
/*         - HTTP Middleware -         */
/* ----------------------------------- */

// GetHTTPMiddlewareAsMuxOptions returns our middleware ready to be passed to the mux.
func GetHTTPMiddlewareAsMuxOptions() []runtime.ServeMuxOption {
	return []runtime.ServeMuxOption{
		runtime.WithErrorHandler(handleHTTPError),
		runtime.WithForwardResponseOption(httpResponseModifier),
	}
}

// handleHTTPError is a custom error handler for the gateway. It's pretty simple.
func handleHTTPError(ctx context.Context, mux *runtime.ServeMux, mar runtime.Marshaler, w http.ResponseWriter, r *http.Request, err error) {
	var (
		grpcStatus                            = status.Convert(err)
		httpStatus, httpResponse, contentType = getHTTPResponseDataFromGRPCStatus(grpcStatus)
		buffer                                = []byte{}
	)

	// 404
	if httpStatus == http.StatusNotFound || httpStatus == http.StatusMethodNotAllowed {
		w.WriteHeader(http.StatusNotFound)
		io.WriteString(w, "not found, check the docs for the correct path and method")
		return
	}

	// 500
	if buffer, err = mar.Marshal(httpResponse); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		io.WriteString(w, "")
		return
	}

	w.Header().Set(contentTypeHeader, contentType)
	w.WriteHeader(httpStatus)
	w.Write(buffer)
}

// getHTTPResponseDataFromGRPCStatus returns the HTTP status, the HTTP response and the content type.
func getHTTPResponseDataFromGRPCStatus(grpcStatus *status.Status) (int, httpErrorResponse, string) {
	httpStatus := runtime.HTTPStatusFromCode(grpcStatus.Code())
	httpResponse := httpErrorResponse{Error: grpcStatus.Message()}
	return httpStatus, httpResponse, contentTypeJSON
}

// httpResponseModifier executes before the response is written to the client.
func httpResponseModifier(ctx context.Context, rw http.ResponseWriter, resp protoreflect.ProtoMessage) error {
	// Delete gRPC-related headers:
	rw.Header().Del("Grpc-Metadata-Content-Type")

	// Add security-related headers:
	rw.Header().Set("X-Content-Type-Options", "nosniff")
	rw.Header().Set("X-Frame-Options", "SAMEORIGIN")
	rw.Header().Set("Strict-Transport-Security", "max-age=31536000; includeSubDomains; preload")
	rw.Header().Set("X-XSS-Protection", "1; mode=block")
	rw.Header().Set("Content-Security-Policy", "default-src 'self'")

	return nil
}
