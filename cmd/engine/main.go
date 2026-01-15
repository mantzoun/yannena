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

	for {
		message, messageError = server.MessagePop()

		if messageError != nil {
			time.Sleep(10 * time.Second)
		} else {
			myLogger.Debug(message)
		}

	}
}
