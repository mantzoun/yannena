package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/websocket"
	"github.com/mantzoun/yannena/internal/logger"
	"github.com/mantzoun/yannena/internal/socket"
)

var (
	myLogger = logger.NewLogger("webui")
)

func main() {
	hub := newHub()
	go hub.run()

	http.HandleFunc("/", serveIndex)
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		serveWS(hub, w, r)
	})
	http.HandleFunc("/command", func(w http.ResponseWriter, r *http.Request) {
		serveCommand(hub, w, r)
	})

	engineClient := socket.NewClient("127.0.0.1:34001")
	engineClient.Connect()
	go getEngineMessages(engineClient)

	http.ListenAndServe(":8080", nil)

	myLogger.Error("Failed to start server")
	os.Exit(1)
}

func getEngineMessages(c *socket.Client) {
	for {
		msg, err := c.MessageWait()
		if err != nil {
			myLogger.Debug("Connection to engine broken")
			time.Sleep(100 * time.Millisecond)
			continue
		}
		c.MessageSend("tok")
		handleEngineMessage(msg)
	}
}

func handleEngineMessage(msg string) {
	myLogger.Debug("received from engine : " + msg)
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

func serveWS(hub *Hub, w http.ResponseWriter, r *http.Request) {
	token := r.URL.Query().Get("token")
	if token == "" {
		http.Error(w, "missing token", http.StatusUnauthorized)
		return
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}

	client := &Client{
		token: token,
		conn:  conn,
		send:  make(chan []byte, 256),
	}

	hub.register <- client

	go writePump(client)
	go readPump(hub, client)
}

func readPump(hub *Hub, c *Client) {
	defer func() {
		hub.unregister <- c
		c.conn.Close()
	}()

	for {
		_, _, err := c.conn.ReadMessage()
		if err != nil {
			return
		}
		// you can ignore client WS input if you want
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

func serveIndex(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "www/index.html")
}

func serveCommand(hub *Hub, w http.ResponseWriter, r *http.Request) {
	token := r.Header.Get("X-User-Token")
	if token == "" {
		http.Error(w, "missing token", http.StatusUnauthorized)
		return
	}

	body, _ := io.ReadAll(r.Body)
	cmd := string(body)

	// Process command in Go code
	response := fmt.Sprintf("OK: %s", cmd)

	// Send private response
	msg := map[string]string{
		"panel": "response",
		"text":  response,
	}
	data, _ := json.Marshal(msg)

	hub.broadcast <- data

	w.Write([]byte(response))
}

//broadcast
//hub.broadcast <- []byte(`{
//  "panel": "updates",
//  "text": "Server tick..."
//}`)
