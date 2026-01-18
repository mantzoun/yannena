package engine

import (
	"encoding/json"
	"os"
	"strconv"

	"github.com/mantzoun/yannena/internal/config"
	"github.com/mantzoun/yannena/internal/logger"
	"github.com/mantzoun/yannena/internal/universe"
)

type Engine struct {
	myLogger *logger.Logger
	saveFile string
	Tik      uint64
	Universe universe.Universe
	config   *config.Config
	//TimeManager timeManager;
	//Utils
}

func NewEngine(config *config.Config) *Engine {
	return &Engine{
		saveFile: config.SaveFile,
		myLogger: logger.NewLogger("engine"),
		config:   config,
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

//Universe universe = Universe(&logger, &timeManager, &utils, utils.idGet(), "Euclid");

//void messageUI(std::string message);

func (e *Engine) TikAdvance() {
	e.Tik += 1
	e.myLogger.Debug("Next turn, " + strconv.FormatUint(e.Tik, 10))

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
