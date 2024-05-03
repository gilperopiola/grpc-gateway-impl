package service

import (
	"context"

	"github.com/gilperopiola/grpc-gateway-impl/app/core/pbs"
)

/* -~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~- */
/*          - Groups Service -         */
/* -~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~- */

func (s *service) CreateGroup(ctx context.Context, req *pbs.CreateGroupRequest) (*pbs.CreateGroupResponse, error) {
	return nil, nil
}

func (s *service) GetGroup(ctx context.Context, req *pbs.GetGroupRequest) (*pbs.GetGroupResponse, error) {
	return nil, nil
}
