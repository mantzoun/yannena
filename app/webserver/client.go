package main

import "github.com/gorilla/websocket"

type Client struct {
	token string
	conn  *websocket.Conn
	send  chan []byte
}

func readPump(hub *Hub, c *Client) {
	defer func() {
		hub.unregister <- c
		c.conn.Close()
	}()

	for {
		myLogger.Debug("loop")
		_, _, err := c.conn.ReadMessage()
		if err != nil {
			return
		}
	}
}

func writePump(c *Client) {
	defer c.conn.Close()

	for msg := range c.send {
		err := c.conn.WriteMessage(websocket.TextMessage, msg)
		if err != nil {
			return
		}
	}
}
