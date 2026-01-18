package main

import (
	"time"

	"github.com/mantzoun/yannena/internal/config"
	"github.com/mantzoun/yannena/internal/engine"
	"github.com/mantzoun/yannena/internal/logger"
	"github.com/mantzoun/yannena/internal/socket"
)

var (
	myLogger   = logger.NewLogger("main")
	configFile = "config.json"
	engineFile = "engine.json"
)

func main() {
	config, _ := config.ReadConfig(configFile)
	var (
		message      string
		messageError error
	)

	myLogger.Debug("Application started")

	server := socket.NewServer(&config)
	server.Start()

	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()

	engine := engine.NewEngine(&config)
	engine.Init()

	for {
		for {
			message, messageError = server.MessagePop()

			if messageError != nil {
				break
			} else {
				myLogger.Debug(message)
			}
		}

		myLogger.Debug("Loop log")
		server.MessageSend("tik")

		engine.TikAdvance()

		<-ticker.C
	}
}
