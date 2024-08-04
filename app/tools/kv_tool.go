package tools

import (
	"context"
	"errors"
)

type CtxKVTool struct {
	context.Context
}

func NewCtxKVTool(ctx context.Context) *CtxKVTool {
	return &CtxKVTool{ctx}
}

func (this *CtxKVTool) AddUserInfoToInnerCtx(userID, username string) {
	this.Context = context.WithValue(this.Context, "userID", userID)
	this.Context = context.WithValue(this.Context, "username", username)
}

func (this *CtxKVTool) GetFromInnerCtx(key string) (string, error) {
	value := this.Context.Value(key)
	if value == nil {
		return "", errors.New("context key not found")
	}
	return value.(string), nil
}

func (this *CtxKVTool) Get(key string) (string, error) {
	return this.GetFromInnerCtx(key)
}

func (this *CtxKVTool) Set(key, value string) error {
	this.Context = context.WithValue(this.Context, key, value) // staticcheck:ignore SA1029
	return nil
}
