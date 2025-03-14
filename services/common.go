package services

import (
	"context"

	"github.com/gin-gonic/gin"
)

type CtxKey string

const (
	CtxKey_Session CtxKey = "CtxKey_Session"
	CtxKey_AppKey  CtxKey = "CtxKey_AppKey"
)

func ToCtx(ctx *gin.Context) context.Context {
	innerCtx := context.Background()
	innerCtx = context.WithValue(innerCtx, CtxKey_AppKey, ctx.GetString(string(CtxKey_AppKey)))
	innerCtx = context.WithValue(innerCtx, CtxKey_Session, ctx.GetString(string(CtxKey_Session)))
	return innerCtx
}

func GetAppKeyFromCtx(ctx context.Context) string {
	if appKey, ok := ctx.Value(CtxKey_AppKey).(string); ok {
		return appKey
	}
	return ""
}
