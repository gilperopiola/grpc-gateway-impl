package app

import (
	"github.com/gilperopiola/grpc-gateway-impl/app/core"
	"github.com/gilperopiola/grpc-gateway-impl/app/core/pbs"
	"github.com/gilperopiola/grpc-gateway-impl/app/core/types/models"
	"github.com/gilperopiola/grpc-gateway-impl/app/servers"
	"github.com/gilperopiola/grpc-gateway-impl/app/service"
	"github.com/gilperopiola/grpc-gateway-impl/app/tools"
)

// This file displays some of the interfaces and structs that are used in the app.
// It doesn't do anything.

/* -~-~-~-~ App Map ~-~-~-~- */

var _ App

// Configuration
var _ core.Config

// Servers
var _ servers.Servers

// Service
var _ service.Service

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

// -> Each one of our Services defined in the protofiles.
var (
	_ core.AuthSvc
	_ core.UsersSvc
	_ core.GroupsSvc
	_ core.HealthSvc

	_ pbs.UnimplementedAuthServiceServer
	_ pbs.UnimplementedUsersServiceServer
	_ pbs.UnimplementedGroupsServiceServer
	_ pbs.UnimplementedHealthServiceServer
)

// -> Tool Interfaces.
// These are all of the actions that any Tools holder can perform.
// Concrete implementations are in the tools pkg.
var (
	_ core.APIs
	_ core.AnyDB
	_ core.DBTool
	_ core.TLSTool
	_ core.CtxTool
	_ core.FileManager
	_ core.ModelConverter
	_ core.PwdHasher
	_ core.RateLimiter
	_ core.RequestValidator
	_ core.ShutdownJanitor
	_ core.TokenGenerator
	_ core.TokenValidator
)

/* -~-~-~-~ Request's Flow ~-~-~-~- */

// When a GRPC Request arrives, our GRPC Server sends it through GRPC Interceptors, and then through the Service
// (which is made of all our different SubServices).
// So: GRPC Server -> Interceptors -> Service.
//
// Our Service, assisted by our Tools (TokenGenerator, PwdHasher, etc), performs certain actions
// (like GetUser or GenerateToken) to get stuff done. These actions sometimes let us communicate with external things,
// like a Database or the File System.
//
// To sum it all up:
// * GRPC Server -> Interceptors -> Service -> Tools -> External Resources (SQL Database, File System, etc).
//
// Oh, and there's also an HTTP Server, but it just adds some middleware and then sends the request through GRPC.
