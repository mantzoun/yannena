package main

import (
	"os"
	"os/signal"
	"syscall"
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
	sigs       = make(chan os.Signal, 1)
	terminate  = false
)

func SignalHandler() {
	sig := <-sigs
	myLogger.Info("Received signal " + sig.String())

	if sig == syscall.SIGUSR1 {
		myLogger.Info("Prepare for exit")
		terminate = true
	}
}

func InitSignalHandler() {
	signal.Notify(sigs, syscall.SIGUSR1)

	go SignalHandler()
}

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

	InitSignalHandler()

	for {
		if terminate {
			myLogger.Info("Exiting")
			os.Exit(0)
		}

		select {
		case <-ticker.C:
			myLogger.Debug("Loop log")
			server.MessageSend("tik")
			engine.TikAdvance()
			engine.Save()
		default:
			for {
				message, messageError = server.MessagePop()
				if messageError != nil {
					break
				} else {
					myLogger.Debug(message)
				}
			}
		}
	}
}
