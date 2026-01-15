package socket

import (
	"errors"

	"github.com/mantzoun/yannena/internal/logger"
)

var (
	mylogger        = logger.NewLogger("socket")
	ErrNoMessages   = errors.New("no messages in queue")
	ErrNoConnection = errors.New("not connected")
)
