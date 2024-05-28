package toolbox

import (
	"context"
	"fmt"

	"github.com/gilperopiola/grpc-gateway-impl/app/core"

	"google.golang.org/grpc/metadata"
)

var _ core.CtxManager = (*ctxManager)(nil)

type ctxManager struct{}

func NewCtxManager() core.CtxManager {
	return &ctxManager{}
}

/* -~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~- */

func (cm ctxManager) AddUserInfo(ctx core.Ctx, userID, username string) core.Ctx {
	ctx = context.WithValue(ctx, &CtxKeyUserID{}, userID)
	ctx = context.WithValue(ctx, &CtxKeyUsername{}, username)
	return ctx
}

func (cm ctxManager) ExtractMetadata(ctx core.Ctx, key string) (string, error) {
	if val := metadata.ValueFromIncomingContext(ctx, key); len(val) > 0 {
		return val[0], nil
	}
	return "", fmt.Errorf("metadata with key %s not found", key)
}

/* -~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~- */

type (
	CtxKeyUserID   struct{}
	CtxKeyUsername struct{}
)
