package breeds

import "github.com/japhy-tech/backend-test/internal/domain/values"

type Breed struct {
	name                values.BreedName
	species             values.Species
	petSize             values.PetSize
	averageFemaleWeight int
	averageMaleWeight   int
}

func (b Breed) Name() values.BreedName {
	return b.name
}

func (b Breed) Species() values.Species {
	return b.species
}

func (b Breed) PetSize() values.PetSize {
	return b.petSize
}

func (b Breed) AverageFemaleWeight() int {
	return b.averageFemaleWeight
}

func (b Breed) AverageMaleWeight() int {
	return b.averageMaleWeight
}
