package socket

import (
	"net"

	"github.com/mantzoun/yannena/internal/config"
)

type Server struct {
	Messenger
	serverStarted bool
	address       string
	ln            net.Listener
	stop          bool
}

func NewServer(config *config.Config) *Server {
	return &Server{
		Messenger: NewMessenger(),
		address:   config.SocketAddress,
	}
}

func (s *Server) Start() {
	s.serverStarted = true
	go serverLoop(s)
}

func (s *Server) Started() bool {
	return s.serverStarted
}

func serverLoop(s *Server) {
	var err error
	s.ln, err = net.Listen("tcp", s.address)

	myLogger.Debug("Start server loop")

	if err != nil {
		myLogger.Error(err.Error())
		s.serverStarted = false
		return
	}

	for {
		if s.stop {
			// for testing, stop the goroutine
			s.ln.Close()
			s.serverStarted = false
			return
		}

		s.conn, err = s.ln.Accept()
		if err != nil {
			myLogger.Debug("Server exiting")
			s.serverStarted = false
			return
		}

		// We only accept one client, modify accordingly for more
		s.ln.Close()

		myLogger.Debug("Client connected. Waiting for messages...")
		s.handleConnection()
		myLogger.Debug("Disconnected from client")

		// Start listening again
		s.ln, err = net.Listen("tcp", s.address)
		if err != nil {
			myLogger.Error("Failed to restart listener: " + err.Error())
			s.serverStarted = false
			return
		}
	}
}

func (s *Server) Stop() {
	s.stop = true
}
