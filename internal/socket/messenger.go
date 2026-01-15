package socket

import (
	"bufio"
	"net"
)

type Messenger struct {
	messages  chan string
	conn      net.Conn
	connected bool
}

func NewMessenger() Messenger {
	return Messenger{}
}

func (m *Messenger) Connected() (bool, error) {
	var err error
	if m.connected {
		err = nil
	} else {
		err = ErrNoConnection
	}
	return m.connected, err
}

func (m *Messenger) Disconnect() {
	if m.conn != nil {
		m.conn.Close()
		m.connected = false
	}
}

func (m *Messenger) MessagePop() (string, error) {
	select {
	case msg := <-m.messages:
		return msg, nil
	default:
		return "", ErrNoMessages
	}
}

func (m *Messenger) MessageSend(message string) (int, error) {
	if m.conn == nil {
		return -1, ErrNoConnection
	}
	return m.conn.Write([]byte(message + "\n"))
}

func (m *Messenger) handleConnection() {
	defer m.conn.Close()
	m.messages = make(chan string, 100)
	scanner := bufio.NewScanner(m.conn)
	for scanner.Scan() {
		message := scanner.Text()

		if message == "" {
			continue
		}

		mylogger.Debug("Received message: " + message)
		m.messages <- message
	}
}
