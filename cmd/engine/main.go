package main

import (
	"time"

	"github.com/mantzoun/yannena/internal/area"
	"github.com/mantzoun/yannena/internal/logger"
	"github.com/mantzoun/yannena/internal/socket"
)

var (
	myLogger   = logger.NewLogger("main")
	areaLogger = logger.NewLogger("area")
)

func main() {
	var (
		message      string
		messageError error
	)

	myLogger.Debug("Application started")

	server := socket.NewServer("127.0.0.1:34001")
	server.Start()

	var test = area.NewBaseArea(areaLogger, 0, "earth")

	areaLogger.Debug(test.Name())

	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()

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

		<-ticker.C
	}
}
