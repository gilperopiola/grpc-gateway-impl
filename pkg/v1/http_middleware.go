package v1

import (
	"context"
	"io"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/reflect/protoreflect"
)

// httpErrorResponse is the struct that gets marshalled onto the HTTP Response when an error occurs.
type httpErrorResponse struct {
	Error string `json:"error"`
}

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
		grpcStatus   = status.Convert(err)
		httpStatus   = runtime.HTTPStatusFromCode(grpcStatus.Code())
		httpResponse = httpErrorResponse{Error: grpcStatus.Message()}
		contentType  = mar.ContentType(grpcStatus)
		buffer       = []byte{}
	)

	if buffer, err = mar.Marshal(httpResponse); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		io.WriteString(w, "")
		return
	}

	w.Header().Set("Content-Type", contentType)
	w.WriteHeader(httpStatus)
	w.Write(buffer)
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
	rw.Header().Set("X-Content-Type-Options", "nosniff")
	rw.Header().Set("Content-Security-Policy", "default-src 'self'")

	return nil
}
