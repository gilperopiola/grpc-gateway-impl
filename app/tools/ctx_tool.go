package tools

import (
	"context"
	"fmt"

	"github.com/gilperopiola/god"
	"github.com/gilperopiola/grpc-gateway-impl/app/core"
	"github.com/gilperopiola/grpc-gateway-impl/app/core/shared/logs"
	"google.golang.org/grpc/metadata"
)

var _ core.ContextManager = ctxTool{}

/* -~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~- */
/* 		    - Context Tool -           */
/* -~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~- */

type ctxTool struct{}

func NewCtxTool() core.ContextManager {
	return &ctxTool{}
}

/* -~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~- */

func (ct ctxTool) AddToCtx(ctx god.Ctx, key, value string) god.Ctx {
	return context.WithValue(ctx, key, value)
}

func (ct ctxTool) GetFromCtx(ctx god.Ctx, key string) (string, error) {
	if value := ctx.Value(key); value != nil {
		return value.(string), nil
	}
	return "", fmt.Errorf("ctx value with key '%s' not found", key)
}

func (ctxTool) GetFromCtxMD(ctx god.Ctx, key string) (string, error) {
	if val := metadata.ValueFromIncomingContext(ctx, key); len(val) > 0 {
		return val[0], nil
	}
	return "", fmt.Errorf("ctx metadata with key %s not found", key)
}

/* -~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~- */

func (ct ctxTool) AddUserInfoToCtx(ctx god.Ctx, userID, username string) god.Ctx {
	ctx = ct.AddToCtx(ctx, CtxKeyUserID, userID)
	ctx = ct.AddToCtx(ctx, CtxKeyUsername, username)
	return ctx
}

// Returns an empty string if there is no user ID in the context.
func (ct ctxTool) GetUserIDFromCtx(ctx god.Ctx) string {
	userID, err := ct.GetFromCtx(ctx, CtxKeyUserID)
	if err != nil {
		logs.LogStrange("Could not get user ID from context: %v", err)
	}
	return userID
}

// Returns an empty string if there is no username in the context.
func (ct ctxTool) GetUsernameFromCtx(ctx god.Ctx) string {
	username, err := ct.GetFromCtx(ctx, CtxKeyUsername)
	if err != nil {
		logs.LogStrange("Could not get username from context: %v", err)
	}
	return username
}

/* -~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~- */

// I know, keys should be struct types.
// But headers come as strings, and I'd rather have it all the same way.
const (
	CtxKeyUserID   = "CtxKeyUserID"
	CtxKeyUsername = "CtxKeyUsername"
)
