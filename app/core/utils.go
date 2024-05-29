package core

import (
	"net/http"

	"github.com/gilperopiola/grpc-gateway-impl/app/core/pbs"
)

// -> Yes, we have a utils file. List of reasons why we shouldn't have one:
// -

// ...
// Get it? It's empty. Utils are fine.

type (
	AuthSvc   = pbs.AuthServiceServer
	UsersSvc  = pbs.UsersServiceServer
	GroupsSvc = pbs.GroupsServiceServer
)

/* -~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~- */

var _ http.ResponseWriter = (*CustomResponseWriter)(nil)

type CustomResponseWriter struct {
	http.ResponseWriter
	Status int
}

func (crw *CustomResponseWriter) WriteHeader(statusCode int) {
	crw.Status = statusCode
	crw.ResponseWriter.WriteHeader(statusCode)
}
