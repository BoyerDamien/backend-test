package values

import (
	"errors"
	"strings"
)

type PetSize int

const (
	Small PetSize = iota
	Medium
	Tall
)

var (
	ErrInvalidPetSize = errors.New("pet size must be one of the following values: [small, medium, tall]")
)

func (p PetSize) String() string {
	switch p {
	case Small:
		return "small"
	case Medium:
		return "medium"
	case Tall:
		return "tall"
	default:
		return ""
	}
}

func PetSizeFromString(s string) (PetSize, error) {
	switch strings.ToLower(s) {
	case "small":
		return Small, nil
	case "medium":
		return Medium, nil
	case "tall":
		return Tall, nil
	default:
		return -1, ErrInvalidPetSize
	}
}
