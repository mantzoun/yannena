package socket

import (
	"strings"
	"testing"
	"time"
)

func execWithRetry(retries int, intervalMS time.Duration, fn func() error) error {
	var err error

	for i := 0; i < retries; i++ {
		err = fn()
		if err == nil {
			return nil
		}
		time.Sleep(intervalMS * time.Millisecond)
	}

	return err
}

func sendMessage(m TestMessenger, message string) {
	m.MessageSend(message)
}

func receiveMessage(m TestMessenger) string {
	var received string
	var err error
	err = execWithRetry(10, 100, func() error {
		var err error
		received, err = m.MessagePop()
		return err
	})
	if err != nil {
		return "no message received"
	}

	return received
}

func verifyMessage(t *testing.T, received string, expected string) {
	if received != strings.TrimSuffix(expected, "\n") {
		t.Error(received + " != " + expected)
		t.Fatal("Message verification failed")
	}
}

func waitForConnection(c *Client) {
	for {
		status, _ := c.Connected()
		if status == true {
			break
		}
		time.Sleep(50 * time.Millisecond)
	}
}

func waitForServerStopped(s *Server) {
	for {
		status := s.Started()
		if status == false {
			break
		}
		time.Sleep(50 * time.Millisecond)
	}
}

func waitForClientStopped(c *Client) {
	for {
		status := c.Started()
		if status == false {
			break
		}
		time.Sleep(50 * time.Millisecond)
	}
}
