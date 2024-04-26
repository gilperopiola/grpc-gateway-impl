package servers

import (
	"context"
	"net/http"
	"strconv"

	"github.com/gilperopiola/grpc-gateway-impl/app/core"
	"github.com/gilperopiola/grpc-gateway-impl/app/core/errs"
	"github.com/gilperopiola/grpc-gateway-impl/app/core/pbs"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
)

/* -~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~- */
/*          - HTTP Gateway -           */
/* -~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~- */

func newHTTPGateway(muxOpts []runtime.ServeMuxOption, middleware middlewareFunc, grpcDialOpts []grpc.DialOption) *http.Server {
	mux := runtime.NewServeMux(muxOpts...)

	ctx := context.Background()
	core.LogPanicIfErr(pbs.RegisterUsersServiceHandlerFromEndpoint(ctx, mux, core.GRPCPort, grpcDialOpts))

	return &http.Server{
		Addr:    core.HTTPPort,
		Handler: middleware(mux),
	}
}

/* -~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~- */
/*         - HTTP Mux Options -        */
/* -~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~- */

// Returns a slice of runtime.ServeMuxOption.
// For now there's only an error handler in here.
func defaultHTTPMuxOpts() []runtime.ServeMuxOption {
	return []runtime.ServeMuxOption{runtime.WithErrorHandler(handleHTTPError)}
}

// handleHTTPError is a custom error handler for the HTTP Gateway. It's pretty simple.
// It converts the GRPC error to an HTTP error and writes it to the response.
func handleHTTPError(ctx context.Context, mux *runtime.ServeMux, m runtime.Marshaler, w http.ResponseWriter, _ *http.Request, err error) {
	grpcStatus := status.Convert(err)
	httpStatus := runtime.HTTPStatusFromCode(grpcStatus.Code())

	// This fn stops the execution chain, so we manually call the forwardResponseOptions to set the headers.
	for _, forwardResponseFn := range mux.GetForwardResponseOptions() {
		forwardResponseFn(ctx, w, nil)
	}

	// Create and marshal an httpError into a []byte buffer. If it fails (unlikely), we return 500.
	var httpBody []byte
	if httpBody, err = m.Marshal(httpError{grpcStatus.Message()}); err != nil {
		core.LogUnexpectedErr(err)
		httpStatus = http.StatusInternalServerError
	}

	writeHTTPErrorResponse(w, httpStatus, httpBody)
}

// Writes the response that gets sent to the client when an error happens.
// For some errors we replace the body with a generic message.
func writeHTTPErrorResponse(rw http.ResponseWriter, status int, body []byte) {
	switch status {
	case http.StatusUnauthorized:
		rw.Header().Set("WWW-Authenticate", "Bearer")
		body = []byte(errs.HTTPUnauthorized)

	case http.StatusForbidden:
		body = []byte(errs.HTTPForbidden)

	case http.StatusNotFound, http.StatusMethodNotAllowed:
		body = []byte(errs.HTTPNotFound)

	case http.StatusInternalServerError:
		body = []byte(errs.HTTPInternal)

	case http.StatusServiceUnavailable:
		body = []byte(errs.HTTPUnavailable)

	default:
		core.LogWeirdBehaviour("HTTP Error Status: " + strconv.Itoa(status))
	}

	rw.WriteHeader(status)
	rw.Write(body)
}

// httpError is the struct that gets marshalled onto the HTTP Response body when an error happens.
// This is what the client would see. The format is '{"error": "error message."}'.
// If this format is changed, then the ErrBodies in errs.go should also change.
type httpError struct {
	Error string `json:"error"`
}
