package socket

import (
	"bufio"
	"errors"
	"net"

	"github.com/mantzoun/yannena/internal/logger"
)

var (
	mylogger      = logger.NewLogger("socket")
	ErrNoMessages = errors.New("no messages in queue")
	messages      = make(chan string, 100)
)

func Start(address string) {
	ln, _ := net.Listen("tcp", address)

	go func() {
		for {
			conn, _ := ln.Accept()
			go handleClient(conn, messages)
		}
	}()
}

func MessagePop() (string, error) {
	select {
	case msg := <-messages:
		return msg, nil
	default:
		return "", ErrNoMessages
	}
}

func handleClient(conn net.Conn, msgChan chan string) {
	defer conn.Close()
	scanner := bufio.NewScanner(conn)
	mylogger.Debug("Client connected. Waiting for messages...")
	for scanner.Scan() {
		msgChan <- scanner.Text()
	}
	mylogger.Debug("Client disconnected.")
}
