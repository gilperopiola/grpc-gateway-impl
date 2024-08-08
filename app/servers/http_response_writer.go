package servers

import (
	"net/http"

	"github.com/gilperopiola/grpc-gateway-impl/app/core"
)

var _ core.HTTPRespWriter = &httpRespWriter{}

/* -~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~- */
/*      - HTTP Response Writer -       */
/* -~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~- */

type httpRespWriter struct {
	http.ResponseWriter
	Body   []byte // nil before first write
	Status *int   // nil before first write
}

func newHTTPRespWriter(rw http.ResponseWriter) core.HTTPRespWriter {
	return &httpRespWriter{
		ResponseWriter: rw,
		Body:           nil,
		Status:         nil,
	}
}

func (rw *httpRespWriter) Write(body []byte) (int, error) {
	rw.Body = append(rw.Body, body...)
	return rw.ResponseWriter.Write(body)
}

func (rw *httpRespWriter) WriteHeader(statusCode int) {
	rw.Status = &statusCode
	rw.ResponseWriter.WriteHeader(statusCode)
}

/* -~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~- */

// Returns nil if the response has not been written yet.
func (rw *httpRespWriter) GetWrittenBody() []byte {
	return rw.Body
}

// Returns a default if the response has not been written yet.
func (rw *httpRespWriter) GetWrittenStatus() int {
	if rw.Status == nil {
		return 0
	}
	return *rw.Status
}
