package everglade

import (
	"github.com/Grant-Eckstein/blind"
)

type Everglade struct {
	Blind  blind.Blind
	Paths  []string
	Layers int
}

func New(l int) (*Everglade, error) {
	var e Everglade
	var err error

	// Set layers
	e.Layers = l

	e.Blind, err = blind.New()
	if err != nil {
		return nil, err
	}
	return &e, nil
}
