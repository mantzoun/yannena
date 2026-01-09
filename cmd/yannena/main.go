package main

import (
	"time"

	"github.com/mantzoun/yannena/internal/logger"
	"github.com/mantzoun/yannena/internal/socket"
)

var (
	prefix = "MAIN"
)

func main() {
	var (
		message      string
		messageError error
	)

	logger.Init()
	logger.Debug(prefix, "Application started")

	socket.Start("127.0.0.1:34001")

	for {
		message, messageError = socket.MessagePop()

		if messageError != nil {
			time.Sleep(10 * time.Second)
		} else {
			logger.Debug(prefix, message)
		}

	}
}
