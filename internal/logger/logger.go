package logger

import (
	"log"
	"os"
	"sync"
)

var (
	myLog *log.Logger
	once  sync.Once
)

func Init() {
	once.Do(func() {
		myLog = log.New(os.Stdout, "", log.LstdFlags)
	})
}

func Debug(prefix string, message string) {
	myLog.Printf("%-10s %-8s %s\n", prefix, "DEBUG", message)
}
