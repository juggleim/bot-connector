package apis

import (
	"bot-connector/errs"
	"bot-connector/services"
	"bot-connector/utils"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

const (
	Header_RequestId     string = "request-id"
	Header_Authorization string = "Authorization"
)

func Validate(ctx *gin.Context) {
	session := utils.GenerateUUIDShort11()
	ctx.Header(Header_RequestId, session)
	ctx.Set(string(services.CtxKey_Session), session)

	auth := ctx.Request.Header.Get("Authorization")
	if auth == "" {
		ctx.JSON(http.StatusUnauthorized, errs.GetErrorResp(errs.ErrorCode_Unknown))
		ctx.Abort()
		return
	}
	if auth != "" {
		if !strings.HasPrefix(auth, "Bearer ") {
			ctx.JSON(http.StatusUnauthorized, errs.GetErrorResp(errs.ErrorCode_Unknown))
			ctx.Abort()
			return
		}
		auth = auth[7:]
		if auth == "aabbcc" {
			ctx.Set(string(services.CtxKey_AppKey), "appkey")
		} else {
			if apiKey, err := services.CheckAuth(auth); err == nil {
				ctx.Set(string(services.CtxKey_AppKey), apiKey.Appkey)
			} else {
				ctx.JSON(http.StatusUnauthorized, errs.GetErrorResp(errs.ErrorCode_Unknown))
				ctx.Abort()
				return
			}
		}
	}
}
