package tools

import (
	"context"
	"fmt"

	"github.com/gilperopiola/god"
	"github.com/gilperopiola/grpc-gateway-impl/app/core"

	"google.golang.org/grpc/metadata"
)

type ctxTool struct{}

func NewCtxTool() core.CtxTool {
	return &ctxTool{}
}

var _ core.CtxTool = ctxTool{}

/* -~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~- */

func (cm ctxTool) AddUserInfo(ctx god.Ctx, userID, username string) god.Ctx {
	ctx = context.WithValue(ctx, &CtxKeyUserID{}, userID)
	ctx = context.WithValue(ctx, &CtxKeyUsername{}, username)
	return ctx
}

func (cm ctxTool) GetMetadata(ctx god.Ctx, key string) (string, error) {
	if val := metadata.ValueFromIncomingContext(ctx, key); len(val) > 0 {
		return val[0], nil
	}
	return "", fmt.Errorf("ctx metadata with key %s not found", key)
}

/* -~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~- */

type (
	CtxKeyUserID   struct{}
	CtxKeyUsername struct{}
)
