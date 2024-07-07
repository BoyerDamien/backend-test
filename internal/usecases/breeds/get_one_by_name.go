package breeds

import (
	"context"

	"github.com/japhy-tech/backend-test/internal/domain/breeds"
	"github.com/japhy-tech/backend-test/internal/domain/values"
	"github.com/japhy-tech/backend-test/internal/usecases"
)

type GetOneByName struct {
	usecases.Base
}

func (g GetOneByName) Info() usecases.UseCaseInfo {
	return usecases.UseCaseInfo{
		Action: usecases.ActionRetrieve,
		Name:   usecases.BreedUsecase,
	}
}

func (g GetOneByName) Handle(ctx context.Context, name string) (*breeds.Breed, error) {
	var (
		breedRepo = g.Datastore().Breeds()
		breedName = values.BreedName(name)
	)

	if err := values.Verify(breedName); err != nil {
		return nil, err
	}
	return breedRepo.GetOneByName(ctx, breedName)
}
