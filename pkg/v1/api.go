package v1

import (
	usersPB "github.com/gilperopiola/grpc-gateway-impl/pkg/users"
	v1Service "github.com/gilperopiola/grpc-gateway-impl/pkg/v1/service"
)

/* ----------------------------------- */
/*             - v1 API -              */
/* ----------------------------------- */

// API is our concrete implementation of the gRPC API defined in the .proto files.
// It implements a handler for each API method, connecting it with the Service.
type API struct {
	Service v1Service.ServiceLayer
	usersPB.UnimplementedUsersServiceServer
}

// NewAPI returns a new instance of the API.
func NewAPI(service v1Service.ServiceLayer) *API {
	return &API{Service: service}
}

// httpErrorResponse is the struct that gets marshalled onto the HTTP Response when an error occurs.
type httpErrorResponse struct {
	Error string `json:"error"`
}
