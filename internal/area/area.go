package area

import (
	"fmt"
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
	PopulationBaseGrowth int64 // points per one hundred thousand
	TechLevel            int
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
	a.populationUpdate()
}

// void rollForNewEffect(void);
// void processActiveEffects(void);
func (area *Area) populationUpdate() {
	var delta int64 = 0
	capacity := 100 - (area.Population*100)/area.PopulationMax

	if capacity > 0 {
		delta = (area.PopulationBaseGrowth * area.Population) / (365 * 100000)
		area.MyLogger.Debug(fmt.Sprintf("Pop delta %d", delta))
	}

	area.AdjustPopulation(delta)
}
