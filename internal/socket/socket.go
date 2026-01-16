package socket

import (
	"errors"

	"github.com/mantzoun/yannena/internal/logger"
)

var (
	myLogger            = logger.NewLogger("socket")
	ErrNoMessages       = errors.New("no messages in queue")
	ErrNoConnection     = errors.New("not connected")
	ErrConnectionClosed = errors.New("connection closed")
)
