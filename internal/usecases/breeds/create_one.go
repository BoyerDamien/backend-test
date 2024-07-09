package breeds

import (
	"context"
	"errors"
	"fmt"

	"github.com/japhy-tech/backend-test/internal/domain/breeds"
	"github.com/japhy-tech/backend-test/internal/domainerror"
	"github.com/japhy-tech/backend-test/internal/usecases"
)

type CreateOne struct {
	usecases.Base
}

func (c CreateOne) Info() usecases.UseCaseInfo {
	return usecases.UseCaseInfo{
		Action: usecases.ActionCreate,
		Name:   usecases.BreedUsecase,
	}
}

func (c CreateOne) Handle(ctx context.Context, params breeds.FactoryOpts) (*breeds.Breed, error) {
	var (
		breedRepo = c.Datastore().Breeds()
	)

	b, err := breeds.NewFactory(params).Instantiate()
	if err != nil {
		return nil, err
	}
	_, err = breedRepo.GetOneByName(ctx, b.Name())
	if err == nil {
		return nil, domainerror.WrapError(domainerror.ErrResourceAlreadyExists, fmt.Errorf("breed %s already exists", b.Name()))
	}
	if !errors.Is(err, domainerror.ErrResourceNotFound) {
		return nil, err
	}
	return breedRepo.CreateOne(ctx, b)
}
