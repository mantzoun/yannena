package universe

import (
	"github.com/mantzoun/yannena/internal/config"
	"github.com/mantzoun/yannena/internal/logger"
	"github.com/mantzoun/yannena/internal/system"
)

type Universe struct {
	MyLogger *logger.Logger
	Systems  map[string]*system.System
}

func NewUniverse(logger *logger.Logger) *Universe {
	return &Universe{
		Systems: make(map[string]*system.System),
	}
}

func (u *Universe) SystemAdd(s *system.System) {
	u.Systems[s.Name] = s
}

func (u *Universe) Init(config *config.Config) {
	u.MyLogger = logger.NewLogger("universe")

	for _, system := range u.Systems {
		system.Init(config)
	}
}

func (u *Universe) TikAdvance() {
	u.MyLogger.Debug("Universe Tik")

	for _, system := range u.Systems {
		system.TikAdvance()
	}
}
