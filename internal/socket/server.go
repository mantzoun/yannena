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
		s.conn, err = s.ln.Accept()
		if err != nil {
			myLogger.Debug("Server exiting")
			s.serverStarted = false
			return
		}
		myLogger.Debug("Client connected. Waiting for messages...")
		s.handleConnection()
		myLogger.Debug("Disconnected from client")
	}
}

func (s *Server) Stop() {
	s.ln.Close()
}
