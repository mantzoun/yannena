package area

import (
	"math"

	"github.com/mantzoun/yannena/internal/config"
	"github.com/mantzoun/yannena/internal/logger"
)

type Area struct {
	Id                   int
	Name                 string
	MyLogger             *logger.Logger
	SystemName           string
	AreaType             AreaType
	Population           int64
	PopulationMax        int64 // penalty after reaching this
	PopulationBaseGrowth int64 // points per thousand
}

func NewArea(name string, id int) *Area {
	return &Area{
		Name: name,
		Id:   id,
	}
}

func (a *Area) Init(config *config.Config) {
	a.MyLogger = logger.NewLogger(a.Name)
}

func (area *Area) AdjustPopulation(delta int64) {
	if area.Population+delta < 0 {
		area.Population = 0
		// TODO planet empty alert
	} else if delta > math.MaxInt64-area.Population {
		area.Population = math.MaxInt64
		// TODO planet full alert
	} else {
		area.Population += delta
	}
}

func (a *Area) TikAdvance() {
	a.MyLogger.Debug("System Tik " + a.Name)
}

// void rollForNewEffect(void);
// void processActiveEffects(void);
func (area *Area) populationUpdate() {
	delta := (area.PopulationBaseGrowth * area.Population) / (365 * 1000)

	area.AdjustPopulation(delta)
}
