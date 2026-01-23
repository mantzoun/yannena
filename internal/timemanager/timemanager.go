package timemanager

import (
	"fmt"

	"github.com/mantzoun/yannena/internal/utilities"
)

type TimeManager struct {
	Turn int
}

func (t *TimeManager) TimeAdvance(tiks int) {
	t.Turn += tiks
}

func (t *TimeManager) TimeGet() int {
	return t.Turn
}

func (t *TimeManager) TimeToString(time int) string {
	if time == 0 {
		return timeToString(t.Turn)
	} else {
		return timeToString(time)
	}
}

func timeToString(time int) string {
	//TODO
	return fmt.Sprintf("Turn %d %d", time, utilities.Roll(0, 100))
}
