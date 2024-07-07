package values

import (
	"errors"
	"strings"
)

type Species int

const (
	Cat Species = iota
	Dog
)

var (
	ErrInvalidSpecies = errors.New("species must be one of the followin values: [dog, cat]")
)

func (s Species) String() string {
	switch s {
	case Cat:
		return "cat"
	case Dog:
		return "dog"
	default:
		return ""
	}
}

func SpeciesFromString(s string) (Species, error) {
	switch strings.ToLower(s) {
	case "cat":
		return Cat, nil
	case "dog":
		return Dog, nil
	default:
		return -1, ErrInvalidSpecies
	}
}
