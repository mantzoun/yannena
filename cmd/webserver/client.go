package main

import "github.com/gorilla/websocket"

type Client struct {
	token string
	conn  *websocket.Conn
	send  chan []byte
}
