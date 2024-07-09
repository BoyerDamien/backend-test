package breeds

import (
	"context"

	"github.com/japhy-tech/backend-test/internal/domain/values"
)

type ListOpts struct {
	Species             *values.Species
	AverageFemaleWeight *int
	AverageMaleWeight   *int
	PetSize             *values.PetSize
	NameIn              []string
}

type Repository interface {
	GetOneByName(context.Context, values.BreedName) (*Breed, error)
	CreateOne(context.Context, *Breed) (*Breed, error)
	UpdateOne(context.Context, *Breed) (*Breed, error)
	DeleteOneByName(context.Context, values.BreedName) error
	List(context.Context, ListOpts) ([]*Breed, error)
	CreateSeveral(context.Context, []*Breed) ([]*Breed, error)
}
