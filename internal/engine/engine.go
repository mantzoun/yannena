package engine

import (
	"encoding/json"
	"os"

	"github.com/mantzoun/yannena/internal/config"
	"github.com/mantzoun/yannena/internal/logger"
	"github.com/mantzoun/yannena/internal/timemanager"
	"github.com/mantzoun/yannena/internal/universe"
)

type Engine struct {
	myLogger *logger.Logger
	saveFile string
	Universe universe.Universe
	config   *config.Config
	TimeMgr  timemanager.TimeManager
	LastId   int
}

func NewEngine(config *config.Config) *Engine {
	return &Engine{
		saveFile: config.SaveFile,
		myLogger: logger.NewLogger("engine"),
		config:   config,
		LastId:   0,
	}
}

func (e *Engine) Init() {
	e.readSaveFile()

	e.Universe.Init(e.config)

	// e.Universe.Init()

	// s := system.NewSystem("Sol", 1, e.myLogger)
	// a := area.NewArea("Earth", 2, e.myLogger)
	// s.AreaAdd(a)
	// a = area.NewArea("Mars", 3, e.myLogger)
	// s.AreaAdd(a)
	// e.Universe.SystemAdd(s)

	// s = system.NewSystem("Alpha Centauri", 4, e.myLogger)
	// a = area.NewArea("Home", 5, e.myLogger)
	// s.AreaAdd(a)
	// e.Universe.SystemAdd(s)

	// e.writeSaveFile()
}

func (e *Engine) Save() {
	e.writeSaveFile()
}

func (e *Engine) TikAdvance() {
	e.TimeMgr.TimeAdvance(1)
	e.myLogger.Debug("Next turn, " + e.TimeMgr.TimeToString(0))

	e.Universe.TikAdvance()
}

func (e *Engine) readSaveFile() error {
	data, err := os.ReadFile(e.saveFile)
	if err != nil {
		return err
	}
	return json.Unmarshal(data, e)
}

func (e *Engine) writeSaveFile() error {
	data, err := json.MarshalIndent(e, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(e.saveFile, data, 0644)
}

func (e *Engine) NewId() int {
	e.LastId += 1
	return e.LastId
}
