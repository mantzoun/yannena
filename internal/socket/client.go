package socket

import (
	"net"
	"time"

	"github.com/mantzoun/yannena/internal/config"
)

type Client struct {
	Messenger
	address string
	stop    bool
	started bool
}

func NewClient(config *config.Config) *Client {
	return &Client{
		Messenger: NewMessenger(),
		address:   config.ClientConnectAddress,
	}
}

func (c *Client) Connect() {
	c.stop = false
	c.started = true
	go clientLoop(c)
}

func clientLoop(c *Client) {
	var err error
	var connectionDelay time.Duration = 100

	myLogger.Debug("Start client loop")

	for {
		if c.stop {
			c.Disconnect()
			myLogger.Debug("Client exiting")
			c.started = false
			return
		}

		c.conn, err = net.DialTimeout("tcp", c.address, 10*time.Millisecond)
		if err != nil {
			myLogger.Debug("Could not connect to server")
			time.Sleep(connectionDelay * time.Millisecond)
			if connectionDelay < 1000 {
				connectionDelay += 100
			}
		} else {
			connectionDelay = 100
			myLogger.Debug("Connected to server")
			c.connected = true
			c.handleConnection()
			c.connected = false
			myLogger.Debug("Disconnected from server")
		}
	}
}

func (c *Client) Started() bool {
	return c.started
}

func (c *Client) Stop() {
	myLogger.Debug("Stopping client")
	c.stop = true
}
