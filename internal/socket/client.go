package socket

import (
	"net"
	"time"
)

type Client struct {
	Messenger
	address string
	stop    bool
	started bool
}

func NewClient(address string) *Client {
	return &Client{
		Messenger: NewMessenger(),
		address:   address,
	}
}

func (c *Client) Connect() {
	c.stop = false
	c.started = true
	go clientLoop(c)
}

func clientLoop(c *Client) {
	var err error

	for {
		if c.stop {
			c.Disconnect()
			mylogger.Debug("Client exiting")
			c.started = false
			return
		}

		c.conn, err = net.DialTimeout("tcp", c.address, 10*time.Millisecond)
		if err != nil {
			mylogger.Debug("Could not connect to server")
			time.Sleep(200 * time.Millisecond)
		} else {
			mylogger.Debug("Connected to server")
			c.connected = true
			c.handleConnection()
			c.connected = false
			mylogger.Debug("Disconnected from server")
		}
	}
}

func (c *Client) Started() bool {
	return c.started
}

func (c *Client) Stop() {
	mylogger.Debug("Stopping client")
	c.stop = true
}
