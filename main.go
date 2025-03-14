package main

import (
	"bot-connector/apis"
	"bot-connector/configures"
	"bot-connector/dbs"
	"bot-connector/logs"
	"bot-connector/services"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	if err := configures.InitConfigures(); err != nil {
		fmt.Println("Init Configures failed.", err)
		return
	}
	logs.InitLogs()
	if err := dbs.InitMysql(); err != nil {
		fmt.Println("Init Mysql failed.", err)
		return
	}

	services.InitTeleBots()

	httpServer := gin.Default()
	httpServer.Use(CorsHandler())
	httpServer.Use(apis.Validate)
	group := httpServer.Group("/bot-connector")
	group.POST("/telebot/events", apis.TeleBotEvents)
	group.POST("/telebot/add", apis.TeleBotAdd)
	group.POST("/telebot/del", apis.TeleBotDel)
	fmt.Println("Start Server with port:", configures.Config.Port)
	httpServer.Run(fmt.Sprintf(":%d", configures.Config.Port))
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
