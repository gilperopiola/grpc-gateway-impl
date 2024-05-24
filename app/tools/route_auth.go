package tools

import (
	"strconv"

	"github.com/gilperopiola/grpc-gateway-impl/app/core"
	"github.com/gilperopiola/grpc-gateway-impl/app/core/errs"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type routeAuthenticator struct {
	// This is usually just core.GetAuthForRoute(). May actually just use it directly.
	getRouteAuthFn func(string) core.RouteAuth
}

func NewRouteAuthenticator(getRouteAuth func(string) core.RouteAuth) core.RouteAuthenticator {
	return &routeAuthenticator{getRouteAuth}
}

func (rauth routeAuthenticator) CanAccessRoute(route, userID string, role core.Role, req any) error {
	switch rauth.getRouteAuthFn(route) {

	case core.RouteAuthPublic:
		return nil

	case core.RouteAuthSelf:
		type PBReqWithUserID interface {
			GetUserId() int32
		}
		reqUserID := req.(PBReqWithUserID).GetUserId()
		if userID != strconv.Itoa(int(reqUserID)) {
			return status.Errorf(codes.PermissionDenied, errs.AuthUserIDInvalid)
		}

	case core.RouteAuthAdmin:
		if role != core.AdminRole {
			core.LogPotentialThreat("User " + userID + " tried to access admin route: " + route)
			return status.Errorf(codes.PermissionDenied, errs.AuthRoleInvalid)
		}

	default:
		core.LogWeirdBehaviour("Route unknown: " + route)
		return status.Errorf(codes.Unknown, errs.AuthRouteUnknown)
	}

	return nil
}
