package breeds

import (
	"github.com/japhy-tech/backend-test/internal/domain/values"
	"github.com/japhy-tech/backend-test/internal/domainerror"
)

type FactoryOpts struct {
	Name                string
	Species             string
	PetSize             string
	AverageFemaleWeight *int
	AverageMaleWeight   *int
}

type Factory struct {
	FactoryOpts
	id int
}

func NewFactory(opts FactoryOpts) *Factory {
	return &Factory{
		FactoryOpts: opts,
	}
}

func (f Factory) SetID(id int) *Factory {
	return &Factory{
		FactoryOpts: f.FactoryOpts,
		id:          id,
	}
}

func (f Factory) Instantiate() (*Breed, error) {
	species, err := values.SpeciesFromString(f.Species)
	if err != nil {
		return nil, domainerror.WrapError(domainerror.ErrDomainValidation, values.ErrInvalidSpecies)
	}

	petSize, err := values.PetSizeFromString(f.PetSize)
	if err != nil {
		return nil, domainerror.WrapError(domainerror.ErrDomainValidation, values.ErrInvalidPetSize)
	}

	if err := values.Verify(values.BreedName(f.Name)); err != nil {
		return nil, err
	}

	return &Breed{
		name: values.BreedName(f.Name),
		averageFemaleWeight: func() int {
			if f.AverageFemaleWeight != nil {
				return *f.AverageFemaleWeight
			}
			return 0
		}(),
		averageMaleWeight: func() int {
			if f.AverageMaleWeight != nil {
				return *f.AverageMaleWeight
			}
			return 0
		}(),
		petSize: petSize,
		species: species,
	}, nil
}
