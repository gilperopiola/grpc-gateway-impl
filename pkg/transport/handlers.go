package transport

import (
	"context"
	"fmt"

	"github.com/gilperopiola/grpc-gateway-impl/pkg/entities"
	usersPB "github.com/gilperopiola/grpc-gateway-impl/pkg/users"

	"github.com/bufbuild/protovalidate-go"
	"google.golang.org/protobuf/proto"
)

type AllPBRequests interface {
	proto.Message
	*usersPB.SignupRequest | *usersPB.LoginRequest
}

type AllPBResponses interface {
	proto.Message
	*usersPB.SignupResponse | *usersPB.LoginResponse
}

type AllServiceRequests interface {
	entities.SignupRequest | entities.LoginRequest
}

type AllServiceResponses interface {
	entities.SignupResponse | entities.LoginResponse
}

func handleRequest[PBReq AllPBRequests, PBRes AllPBResponses, SvcReq AllServiceRequests, SvcRes AllServiceResponses](
	ctx context.Context,
	pbRequest PBReq,
	toService func(PBReq) (SvcReq, error),
	serviceMethod func(context.Context, SvcReq) (SvcRes, error),
	toTransport func(SvcRes) (PBRes, error),
	pbValidator *protovalidate.Validator,
) (PBRes, error) {

	// PB Request -> Gets validated
	if err := pbValidator.Validate(pbRequest); err != nil {
		return nil, fmt.Errorf("invalid signup request: %w", err)
	}

	// PB Request -> To Service Request
	serviceRequest, err := toService(pbRequest)
	if err != nil {
		return nil, fmt.Errorf("error converting pb request to service format: %w", err)
	}

	// Call Service Method with Service Request
	serviceResponse, err := serviceMethod(ctx, serviceRequest)
	if err != nil {
		return nil, fmt.Errorf("service error: %w", err)
	}

	// Service Response -> To PB Response
	pbResponse, err := toTransport(serviceResponse)
	if err != nil {
		return nil, fmt.Errorf("error converting service response to transport format: %w", err)
	}

	return pbResponse, nil
}

func (t *transport) Signup(ctx context.Context, pbRequest *usersPB.SignupRequest) (*usersPB.SignupResponse, error) {
	return handleRequest(ctx, pbRequest, t.toService.Signup, t.service.Signup, t.toTransport.Signup, t.protoValidator)
}

func (t *transport) Login(ctx context.Context, pbRequest *usersPB.LoginRequest) (*usersPB.LoginResponse, error) {
	return handleRequest(ctx, pbRequest, t.toService.Login, t.service.Login, t.toTransport.Login, t.protoValidator)
}
