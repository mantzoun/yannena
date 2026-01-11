package area

import (
	"github.com/mantzoun/yannena/internal/entity"
	"github.com/mantzoun/yannena/internal/logger"
)

type AreaType int

const (
	StateIdle      AreaType = iota // 0
	StateConnected                 // 1
	StateError                     // 2
	StateRetrying                  // 3
)

type BaseArea struct {
	*entity.BaseEntity
	systemName string
	areaType   AreaType
	// pop_t population;
	// pop_t populationMax;
	// int   populationBaseGrowth; // points per thousand
}

func NewBaseArea(logger *logger.Logger, id int, name string) *BaseArea {
	return &BaseArea{
		BaseEntity: entity.NewBaseEntity(logger, id, name),
	}
}

func (area *BaseArea) SystemName() string {
	return area.systemName
}

func (area *BaseArea) SetSystemName(name string) {
	if name != "" {
		area.systemName = name
	}
}

func (area *BaseArea) AreaType() AreaType {
	return area.areaType
}

func (area *BaseArea) SetAreaType(areaType AreaType) {
	area.areaType = areaType
}

func (area *BaseArea) tikAdvance() {
	area.MyLogger.Debug("ADVANCE")
}

//   void rollForNewEffect(void);
//   void processActiveEffects(void);
//   void populationUpdate(void);

//   void init();

//   void populationSet(pop_t pop);
//   void populationMod(pop_t diff);
//   pop_t populationGet(void);

//   void populationBaseGrowthSet(int b);
//   int populationBaseGrowthGet(void);
