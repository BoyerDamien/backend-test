package breeds

import (
	"context"
	"errors"

	"github.com/japhy-tech/backend-test/internal/domain/breeds"
	"github.com/japhy-tech/backend-test/internal/domainerror"
	"github.com/japhy-tech/backend-test/internal/usecases"
)

type UpdateOne struct {
	usecases.Base
}

func (c UpdateOne) Info() usecases.UseCaseInfo {
	return usecases.UseCaseInfo{
		Name:   usecases.BreedUsecase,
		Action: usecases.ActionUpdate,
	}
}

func (c UpdateOne) Handle(ctx context.Context, params breeds.FactoryOpts) (*breeds.Breed, error) {
	var (
		breedRepo = c.Datastore().Breeds()
	)

	b, err := breeds.NewFactory(params).Instantiate()
	if err != nil {
		return nil, err
	}
	_, err = breedRepo.GetOneByName(ctx, b.Name())
	if errors.Is(err, domainerror.ErrResourceNotFound) {
		return nil, err
	}
	return breedRepo.UpdateOne(ctx, b)
}
