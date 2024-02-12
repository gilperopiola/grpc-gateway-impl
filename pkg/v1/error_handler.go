package v1

import (
	"context"
	"io"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc/status"
)

/* ----------------------------------- */
/*          - Error Handler -          */
/* ----------------------------------- */

// GetServerOptions returns our error handler ready to be passed to the mux.
func GetServerOptions() runtime.ServeMuxOption {
	return runtime.WithErrorHandler(httpErrorHandler)
}

// httpErrorResponse is the struct that gets marshalled onto the HTTP Response when an error occurs.
type httpErrorResponse struct {
	Error string `json:"error"`
}

// httpErrorHandler is a custom error handler for the gateway. It's pretty simple.
func httpErrorHandler(ctx context.Context, mux *runtime.ServeMux, mar runtime.Marshaler, w http.ResponseWriter, r *http.Request, err error) {
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
