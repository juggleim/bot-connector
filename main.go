package main

import (
	"bot-connector/configures"
	"bot-connector/dbs"
	"bot-connector/logs"
	"bot-connector/routers"
	"bot-connector/services"
	"fmt"

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
	routers.Route(httpServer, "bot-connector")
	fmt.Println("Start Server with port:", configures.Config.Port)
	httpServer.Run(fmt.Sprintf(":%d", configures.Config.Port))
}
