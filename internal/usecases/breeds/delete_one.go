package breeds

import (
	"context"

	"github.com/japhy-tech/backend-test/internal/domain/values"
	"github.com/japhy-tech/backend-test/internal/usecases"
)

type DeleteOneByName struct {
	usecases.Base
}

func (d DeleteOneByName) Info() usecases.UseCaseInfo {
	return usecases.UseCaseInfo{
		Name:   usecases.BreedUsecase,
		Action: usecases.ActionDelete,
	}
}

func (d DeleteOneByName) Handle(ctx context.Context, name string) error {
	var (
		breedRepo = d.Datastore().Breeds()
		breedName = values.BreedName(name)
	)

	if err := values.Verify(breedName); err != nil {
		return err
	}

	_, err := breedRepo.GetOneByName(ctx, breedName)
	if err != nil {
		return err
	}
	return breedRepo.DeleteOneByName(ctx, breedName)
}
