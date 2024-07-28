package tools

import (
	"net/http"

	"github.com/gilperopiola/grpc-gateway-impl/app/core"
)

var _ core.HTTPRespWriter = (*httpRespWriter)(nil)

type httpRespWriter struct {
	http.ResponseWriter

	Body   []byte
	Status int
}

func NewHTTPRespWriter(rw http.ResponseWriter) core.HTTPRespWriter {
	return &httpRespWriter{
		ResponseWriter: rw,
		Body:           []byte{},
		Status:         200, // rw default
	}
}

func (rw *httpRespWriter) Write(p []byte) (int, error) {
	rw.Body = append(rw.Body, p...)
	return rw.ResponseWriter.Write(p)
}

func (rw *httpRespWriter) WriteHeader(statusCode int) {
	rw.Status = statusCode
	rw.ResponseWriter.WriteHeader(statusCode)
}

// Returns []byte{} if response not written yet.
func (rw *httpRespWriter) GetWrittenBody() []byte {
	return rw.Body
}

// Returns 0 if status not written yet.
func (rw *httpRespWriter) GetWrittenStatus() int {
	return rw.Status
}
