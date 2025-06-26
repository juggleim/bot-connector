package apis

import (
	"bot-connector/dbs"
	"bot-connector/errs"
	"bot-connector/services"
	"bot-connector/utils"
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

const (
	Header_RequestId     string = "request-id"
	Header_Authorization string = "Authorization"

	Header_AppKey    string = "appkey"
	Header_Nonce     string = "nonce"
	Header_Timestamp string = "timestamp"
	Header_Signature string = "signature"
)

func Validate(ctx *gin.Context) {
	session := utils.GenerateUUIDShort11()
	ctx.Header(Header_RequestId, session)
	ctx.Set(string(services.CtxKey_Session), session)

	auth := ctx.Request.Header.Get("Authorization")
	if auth != "" {
		if !strings.HasPrefix(auth, "Bearer ") {
			ctx.JSON(http.StatusUnauthorized, errs.GetErrorResp(errs.ErrorCode_Unknown))
			ctx.Abort()
			return
		}
		auth = auth[7:]
		if apiKey, err := services.CheckAuth(auth); err == nil {
			ctx.Set(string(services.CtxKey_AppKey), apiKey.Appkey)
		} else {
			ctx.JSON(http.StatusUnauthorized, errs.GetErrorResp(errs.ErrorCode_Unknown))
			ctx.Abort()
			return
		}
	} else {
		appKey := ctx.Request.Header.Get(Header_AppKey)
		nonce := ctx.Request.Header.Get(Header_Nonce)
		tsStr := ctx.Request.Header.Get(Header_Timestamp)
		signature := ctx.Request.Header.Get(Header_Signature)
		if appKey == "" {
			ctx.JSON(http.StatusBadRequest, errs.GetErrorResp(errs.ErrorCode_APPKEY_REQUIRED))
			ctx.Abort()
			return
		}
		if nonce == "" {
			ctx.JSON(http.StatusBadRequest, errs.GetErrorResp(errs.ErrorCode_NONCE_REQUIRED))
			ctx.Abort()
			return
		}
		if tsStr == "" {
			ctx.JSON(http.StatusBadRequest, errs.GetErrorResp(errs.ErrorCode_TIMESTAMP_REQUIRED))
			ctx.Abort()
			return
		}
		if signature == "" {
			ctx.JSON(http.StatusBadRequest, errs.GetErrorResp(errs.ErrorCode_SIGNATURE_REQUIRED))
			ctx.Abort()
			return
		}
		dao := dbs.AppInfoDao{}
		appinfo := dao.FindByAppkey(appKey)
		if appinfo == nil {
			ctx.JSON(http.StatusBadRequest, errs.GetErrorResp(errs.ErrorCode_APP_NOT_EXISTED))
			ctx.Abort()
			return
		}
		str := fmt.Sprintf("%s%s%s", appinfo.AppSecret, nonce, tsStr)
		sig := SHA1(str)
		if sig == signature {
			ctx.Set(string(services.CtxKey_AppKey), appKey)
		} else {
			ctx.JSON(http.StatusForbidden, errs.ErrorCode_SIGNATURE_FAIL)
			ctx.Abort()
			return
		}
	}
}

func SHA1(s string) string {
	o := sha1.New()
	o.Write([]byte(s))
	return hex.EncodeToString(o.Sum(nil))
}
