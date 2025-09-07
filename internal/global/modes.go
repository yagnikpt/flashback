package global

import (
	"log"
)

type ModeState int

const (
	StateNote ModeState = iota
	StateRecall
	StateDelete
)

var stateName = map[ModeState]string{
	StateNote:   "note",
	StateRecall: "recall",
	StateDelete: "delete",
}

func (ss ModeState) String() string {
	return stateName[ss]
}

func (ss *ModeState) NextMode() {
	switch *ss {
	case StateNote:
		*ss = StateRecall
	case StateRecall:
		*ss = StateDelete
	case StateDelete:
		*ss = StateNote
	default:
		log.Printf("unknown state: %s", ss)
		*ss = StateNote
	}
}
