package main

import (
	"io"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/websocket"
	"github.com/mantzoun/yannena/internal/config"
	"github.com/mantzoun/yannena/internal/logger"
	"github.com/mantzoun/yannena/internal/socket"
)

var (
	myLogger     = logger.NewLogger("webui")
	confFile     = "config.json"
	engineClient *socket.Client
)

func main() {
	config, _ := config.ReadConfig(confFile)
	hub := newHub()
	go hub.run()

	http.HandleFunc("/", serveIndex)
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		serveWS(hub, w, r)
	})
	http.HandleFunc("/command", func(w http.ResponseWriter, r *http.Request) {
		serveCommand(hub, w, r)
	})
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		serveHealth(hub, w, r)
	})

	engineClient = socket.NewClient(&config)
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

func serveIndex(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "www/index.html")
}

func serveHealth(hub *Hub, w http.ResponseWriter, r *http.Request) {
	connected, _ := engineClient.Connected()
	if connected {
		w.Write([]byte("Health OK"))
	} else {
		http.Error(w, "Not connected to engime", http.StatusInternalServerError)
	}
}

func serveCommand(hub *Hub, w http.ResponseWriter, r *http.Request) {
	token := r.Header.Get("X-User-Token")
	if token == "" {
		http.Error(w, "missing token", http.StatusUnauthorized)
		return
	}

	body, _ := io.ReadAll(r.Body)
	cmd := string(body)

	myLogger.Debug(token + " message: " + cmd)
	engineClient.MessageSend(cmd)
	// Process command in Go code
	// response := fmt.Sprintf("OK: %s", cmd)

	// Send private response
	//	msg := map[string]string{
	//		"panel": "response",
	//		"text":  "received OK",
	//	}
	//	data, _ := json.Marshal(msg)

	//	hub.broadcast <- data

	w.Write([]byte("received OK"))
}

//broadcast
//hub.broadcast <- []byte(`{
//  "panel": "updates",
//  "text": "Server tick..."
//}`)
