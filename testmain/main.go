package main

import (
	"bot-connector/configures"
	"bot-connector/services"
	"fmt"
)

func main() {
	if err := configures.InitConfigures(); err != nil {
		fmt.Println("Init Configures failed.", err)
		return
	}
	apiKey, err := services.GenerateApiKey("appkey", "", "")
	fmt.Println(apiKey, err)
}
