package socket

import (
	"net"
)

type Server struct {
	Messenger
	serverStarted bool
	address       string
	ln            net.Listener
}

func NewServer(address string) *Server {
	return &Server{
		Messenger: NewMessenger(),
		address:   address,
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

	if err != nil {
		mylogger.Error(err.Error())
		s.serverStarted = false
		return
	}

	for {
		s.conn, err = s.ln.Accept()
		if err != nil {
			mylogger.Debug("Server exiting")
			s.serverStarted = false
			return
		}
		mylogger.Debug("Client connected. Waiting for messages...")
		s.handleConnection()
		mylogger.Debug("Disconnected from client")
	}
}

func (s *Server) Stop() {
	s.ln.Close()
}
