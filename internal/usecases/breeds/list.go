package breeds

import (
	"context"

	"github.com/japhy-tech/backend-test/internal/domain/breeds"
	"github.com/japhy-tech/backend-test/internal/domain/values"
	"github.com/japhy-tech/backend-test/internal/usecases"
)

type List struct {
	usecases.Base
}

type ListOpts struct {
	Species             *string
	AverageFemaleWeight *int
	AverageMaleWeight   *int
}

func (g List) Info() usecases.UseCaseInfo {
	return usecases.UseCaseInfo{
		Action: usecases.ActionList,
		Name:   usecases.BreedUsecase,
	}
}

func (g List) Handle(ctx context.Context, params ListOpts) ([]*breeds.Breed, error) {
	var opts breeds.ListOpts

	if params.Species != nil {
		if val, err := values.SpeciesFromString(*params.Species); err != nil {
			return nil, err
		} else {
			opts.Species = &val
		}
	}
	opts.AverageMaleWeight = params.AverageMaleWeight
	opts.AverageFemaleWeight = params.AverageFemaleWeight

	return g.Datastore().Breeds().List(ctx, opts)
}
