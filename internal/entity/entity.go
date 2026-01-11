package entity

import "github.com/mantzoun/yannena/internal/logger"

type Entity interface {
	tikAdvance()
}

type BaseEntity struct {
	id       int
	name     string
	MyLogger *logger.Logger
	// timemanager
	// utils
}

func NewBaseEntity(logger *logger.Logger, id int, name string) *BaseEntity {
	return &BaseEntity{
		MyLogger: logger,
		id:       id,
		name:     name,
	}
}

func (entity *BaseEntity) Name() string {
	return entity.name
}

func (entity *BaseEntity) SetName(name string) {
	if name != "" {
		entity.name = name
	}
}

func (entity *BaseEntity) Id() int {
	return entity.id
}

func (entity *BaseEntity) SetId(id int) {
	entity.id = id
}
