package app

import (
	"github.com/gilperopiola/grpc-gateway-impl/app/core"
	"github.com/gilperopiola/grpc-gateway-impl/app/core/errs"
	"github.com/gilperopiola/grpc-gateway-impl/app/core/models"
	"github.com/gilperopiola/grpc-gateway-impl/app/core/pbs"
	"github.com/gilperopiola/grpc-gateway-impl/app/servers"
	"github.com/gilperopiola/grpc-gateway-impl/app/service"
	"github.com/gilperopiola/grpc-gateway-impl/app/tools"
)

// This file has some of the interfaces and structs that are used in the app.
// It doesn't do anything.

/* -~-~-~-~ App Map ~-~-~-~- */

var _ App

// Configuration
var _ core.Config

// Servers
var _ servers.Servers

// Services
var _ service.Services

// Tools
var _ core.Tools
var _ tools.Tools

var (
	// Routes
	_ = core.Routes

	// Auth-Level required per route
	_ = models.RouteAuthPublic
	_ = models.RouteAuthUser
	_ = models.RouteAuthSelf
	_ = models.RouteAuthAdmin

	// Roles
	_ = models.DefaultRole
	_ = models.AdminRole
)

// -> SQL Database Models.
// Used to migrate the DB.
var _ = models.AllDBModels

// -> Each one of our Services as defined in the proto.
var (
	_ core.AuthSvc
	_ core.UsersSvc
	_ core.GroupsSvc

	_ pbs.UnimplementedAuthServiceServer
	_ pbs.UnimplementedUsersServiceServer
	_ pbs.UnimplementedGroupsServiceServer
)

// -> Tools / Tools Interfaces.
// These are all of the Actions that any Tools holder can perform.
// Concrete implementations are in the tools pkg.
var (
	_ core.ExternalAPIs
	_ core.AnyDB
	_ core.DBTool
	_ core.TLSTool
	_ core.CtxTool
	_ core.FileManager
	_ core.ModelConverter
	_ core.PwdHasher
	_ core.RateLimiter
	_ core.RequestsValidator
	_ core.ShutdownJanitor
	_ core.TokenGenerator
	_ core.TokenValidator
)

// DB Layer and Service Layer Errors.
var _ errs.DBErr
var _ errs.ServiceErr
