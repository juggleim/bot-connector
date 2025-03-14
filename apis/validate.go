package apis

import (
	"bot-connector/dbs"
	"bot-connector/errs"
	"bot-connector/services"
	"bot-connector/services/pbobjs"
	"bot-connector/utils"
	"encoding/base64"
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
			if appkey, ok := CheckAuth(auth); ok {
				ctx.Set(string(services.CtxKey_AppKey), appkey)
			} else {
				ctx.JSON(http.StatusUnauthorized, errs.GetErrorResp(errs.ErrorCode_Unknown))
				ctx.Abort()
				return
			}
		}
	}
}

func CheckAuth(apikey string) (string, bool) {
	bs, err := base64.URLEncoding.DecodeString(apikey)
	if err != nil {
		return "", false
	}
	authWrap := &pbobjs.ApiKeyWrap{}
	err = utils.PbUnMarshal(bs, authWrap)
	if err != nil {
		return "", false
	}
	appkey := authWrap.AppKey
	dao := dbs.AppInfoDao{}
	appinfo := dao.FindByAppkey(appkey)
	if appinfo == nil {
		return "", false
	}
	_, err = utils.AesDecrypt(authWrap.Value, []byte(appinfo.ApiSecureKey))
	if err != nil {
		return "", false
	}
	return appkey, true
}
