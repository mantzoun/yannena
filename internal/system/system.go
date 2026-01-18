package system

import (
	"github.com/mantzoun/yannena/internal/area"
	"github.com/mantzoun/yannena/internal/config"
	"github.com/mantzoun/yannena/internal/logger"
)

type System struct {
	Id       int
	Name     string
	MyLogger *logger.Logger
	Areas    map[string]*area.Area
}

func NewSystem(name string, id int, logger *logger.Logger) *System {
	return &System{
		Name:  name,
		Id:    id,
		Areas: make(map[string]*area.Area),
	}
}

func (s *System) Init(config *config.Config) {
	s.MyLogger = logger.NewLogger(s.Name)

	for _, area := range s.Areas {
		area.Init(config)
	}
}

func (s *System) AreaAdd(a *area.Area) {
	s.Areas[a.Name] = a
}

func (s *System) TikAdvance() {
	s.MyLogger.Debug("System Tik " + s.Name)

	for _, area := range s.Areas {
		area.TikAdvance()
	}
}
