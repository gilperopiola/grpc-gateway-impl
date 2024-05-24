package app

import (
	"github.com/gilperopiola/grpc-gateway-impl/app/core"
	"github.com/gilperopiola/grpc-gateway-impl/app/core/errs"
)

var (
	_ App
	_ Tools
	_ core.Config
	_ core.Servers
	_ core.Service
)

var (
	_ core.AuthSvc
	_ core.UsersSvc
	_ core.GroupsSvc
)

var (
	_ core.Toolbox
	_ core.APIs
	_ core.DB
	_ core.DBTool
	_ core.FileManager
	_ core.MetadataGetter
	_ core.PwdHasher
	_ core.RateLimiter
	_ core.RequestsValidator
	_ core.RouteAuthenticator
	_ core.ShutdownJanitor
	_ core.TLSTool
	_ core.TokenGenerator
	_ core.TokenValidator
)

var (
	_ core.GRPCInfo
	_ core.GRPCHandler
	_ core.GRPCInterceptors
	_ core.GRPCServerOptions
	_ core.GRPCDialOptions
	_ core.GRPCServiceRegistrar

	_ core.HTTPMultiplexer

	_ errs.DBErr
	_ errs.ServiceErr
)
