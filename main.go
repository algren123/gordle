package main

import (
	"fmt"

	"github.com/algren123/gordle/bot"
	"github.com/algren123/gordle/config"
)

func main() {
	err := config.ReadConfig()
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	bot.Start()

	<-make(chan struct{})
	return
}
