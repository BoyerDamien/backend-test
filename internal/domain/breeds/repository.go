package breeds

import (
	"context"

	"github.com/japhy-tech/backend-test/internal/domain/values"
)

type ListOpts struct {
	Species             values.Species
	AverageFemaleWeight int
	AverageMaleWeight   int
}

type Repository interface {
	GetOneByName(context.Context, values.BreedName) (*Breed, error)
	Save(context.Context, *Breed) (*Breed, error)
	UpdateOrCreateOneByName(context.Context, values.BreedName, *Breed) (*Breed, error)
	DeleteOneByName(context.Context, values.BreedName) error
	List(context.Context, ListOpts) ([]*Breed, error)
}
