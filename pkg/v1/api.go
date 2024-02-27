package v1

import (
	usersPB "github.com/gilperopiola/grpc-gateway-impl/pkg/users"
	"github.com/gilperopiola/grpc-gateway-impl/pkg/v1/service"
)

/* ----------------------------------- */
/*             - v1 API -              */
/* ----------------------------------- */

// API is our concrete implementation of the gRPC API defined in the .proto files.
// It has a handler for each API method, connecting it with the Service.
// It implements the usersPB.UsersServiceServer interface.
type API struct {
	Service service.ServiceLayer
	usersPB.UnimplementedUsersServiceServer
}

// NewAPI returns a new instance of the API.
func NewAPI(service service.ServiceLayer) *API {
	return &API{Service: service}
}
