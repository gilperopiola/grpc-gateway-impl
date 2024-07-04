package app

import (
	"github.com/gilperopiola/grpc-gateway-impl/app/core"
	"github.com/gilperopiola/grpc-gateway-impl/app/core/errs"
	"github.com/gilperopiola/grpc-gateway-impl/app/core/models"
	"github.com/gilperopiola/grpc-gateway-impl/app/core/pbs"
	"github.com/gilperopiola/grpc-gateway-impl/app/servers"
	"github.com/gilperopiola/grpc-gateway-impl/app/service"
	"github.com/gilperopiola/grpc-gateway-impl/app/toolbox"
)

// -> Our main App.
var _ App

var (
	// -> The components of our App.
	_ core.Config
	_ core.Servers
	_ core.Service
	_ core.Toolbox

	_ servers.Servers
	_ service.Service
	_ toolbox.Toolbox
)

// -> DB Models.
var _ = models.AllDBModels

var (
	// -> Routes and Auth.
	_ = core.Routes
	_ = core.AuthForRoute("with the route name you get the auth required for it")
	_ = core.CanAccessRoute("route", "user_id", "role", "request")

	_ = models.RouteAuthPublic
	_ = models.RouteAuthUser
	_ = models.RouteAuthSelf
	_ = models.RouteAuthAdmin

	_ = models.DefaultRole
	_ = models.AdminRole
)

var (
	// -> Each Service.
	_ core.AuthSvc
	_ core.UsersSvc
	_ core.GroupsSvc

	_ pbs.UnimplementedAuthServiceServer
	_ pbs.UnimplementedUsersServiceServer
	_ pbs.UnimplementedGroupsServiceServer
)

var (
	// -> Tool Interfaces.
	// Concrete implementations are in the toolbox pkg.
	_ core.APIs
	_ core.DB
	_ core.DBTool
	_ core.TLSTool
	_ core.CtxManager
	_ core.FileManager
	_ core.ModelConverter
	_ core.PwdHasher
	_ core.RateLimiter
	_ core.Retrier
	_ core.RequestsValidator
	_ core.ShutdownJanitor
	_ core.TokenGenerator
	_ core.TokenValidator
)

// -> DB Layer and Service Layer Errors.
var _ errs.DBErr
var _ errs.ServiceErr
