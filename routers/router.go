package routers

import (
	"bot-connector/apis"
	"net/http"

	"github.com/gin-gonic/gin"
)

func Route(eng *gin.Engine, prefix string) {
	eng.Use(CorsHandler())
	group := eng.Group("/" + prefix)
	group.Use(apis.Validate)

	group.POST("/telebot/events", apis.TeleBotEvents)
	group.POST("/telebot/add", apis.TeleBotAdd)
	group.POST("/telebot/del", apis.TeleBotDel)
}

func CorsHandler() gin.HandlerFunc {
	return func(context *gin.Context) {
		method := context.Request.Method
		context.Writer.Header().Add("Access-Control-Allow-Origin", "*")
		context.Writer.Header().Add("Access-Control-Allow-Headers", "Content-Type,AccessToken,X-CSRF-Token, Authorization")
		context.Writer.Header().Add("Access-Control-Allow-Methods", "POST, GET, OPTIONS, DELETE, PATCH, PUT")
		context.Writer.Header().Add("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers, Content-Type")
		context.Writer.Header().Add("Access-Control-Allow-Credentials", "true")

		if method == "OPTIONS" {
			context.AbortWithStatus(http.StatusNoContent)
		}
		context.Next()
	}
}
