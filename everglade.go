package everglade

import (
	"github.com/Grant-Eckstein/blind"
)

type Everglade struct {
	Blind blind.Blind
}

func New() (*Everglade, error) {
	var e Everglade
	var err error
	e.Blind, err = blind.New()
	if err != nil {
		return nil, err
	}
	return &e, nil
}
