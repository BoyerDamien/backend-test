package mysql

import (
	"context"
	"fmt"

	"github.com/doug-martin/goqu/v9"
	"github.com/japhy-tech/backend-test/internal/common"
	"github.com/japhy-tech/backend-test/internal/domain/breeds"
	"github.com/japhy-tech/backend-test/internal/domain/values"
	"github.com/japhy-tech/backend-test/internal/domainerror"
)

type BreedStorage struct {
	db *goqu.Database
}

func NewBreedStorage(db *goqu.Database) *BreedStorage {
	return &BreedStorage{
		db: db,
	}
}

type BreedModel struct {
	Name                     string `db:"name"`
	Species                  string `db:"species"`
	PetSize                  string `db:"pet_size"`
	AverageMaleAdultWeight   int    `db:"average_male_adult_weight"`
	AverageFemaleAdultWeight int    `db:"average_female_adult_weight"`
}

func (b BreedModel) ToDomain() (*breeds.Breed, error) {
	return breeds.NewFactory(breeds.FactoryOpts{
		Name:                b.Name,
		Species:             b.Species,
		PetSize:             b.PetSize,
		AverageFemaleWeight: &b.AverageFemaleAdultWeight,
		AverageMaleWeight:   &b.AverageMaleAdultWeight,
	}).Instantiate()
}

func (b BreedStorage) GetOneByName(ctx context.Context, name values.BreedName) (*breeds.Breed, error) {
	var res BreedModel

	query := b.db.From("breeds").
		Select(
			"name",
			"pet_size",
			"average_male_adult_weight",
			"average_female_adult_weight",
			"species",
		).
		Where(goqu.I("name").Eq(name.String()))
	found, err := query.ScanStructContext(ctx, &res)
	if err != nil {
		return nil, domainerror.WrapError(domainerror.ErrInternalError, err)
	}
	if !found {
		return nil, domainerror.WrapError(domainerror.ErrResourceNotFound, fmt.Errorf("breed %s not found", name))
	}

	return res.ToDomain()
}

func (b BreedStorage) CreateOne(ctx context.Context, input *breeds.Breed) (*breeds.Breed, error) {
	insert := b.db.Insert(goqu.T("breeds")).Rows(
		goqu.Record{
			"name":                        input.Name().String(),
			"species":                     input.Species().String(),
			"pet_size":                    input.PetSize().String(),
			"average_male_adult_weight":   input.AverageMaleWeight(),
			"average_female_adult_weight": input.AverageFemaleWeight(),
		},
	).Executor()

	if _, err := insert.ExecContext(ctx); err != nil {
		return nil, domainerror.WrapError(domainerror.ErrInternalError, err)
	}
	return b.GetOneByName(ctx, input.Name())
}

func (b BreedStorage) UpdateOne(ctx context.Context, input *breeds.Breed) (*breeds.Breed, error) {
	update := b.db.Update(goqu.T("breeds")).
		Set(goqu.Record{
			"name":                        input.Name().String(),
			"species":                     input.Species().String(),
			"pet_size":                    input.PetSize().String(),
			"average_male_adult_weight":   input.AverageMaleWeight(),
			"average_female_adult_weight": input.AverageFemaleWeight(),
		}).
		Where(
			goqu.C("name").Eq(input.Name()),
		).Executor()

	res, err := update.ExecContext(ctx)
	if err != nil {
		return nil, domainerror.WrapError(domainerror.ErrInternalError, err)
	}
	if i, err := res.RowsAffected(); err != nil {
		return nil, domainerror.WrapError(domainerror.ErrInternalError, err)
	} else if i == 0 {
		return nil, domainerror.ErrNothingTodo
	}

	return b.GetOneByName(ctx, input.Name())
}

func (b BreedStorage) DeleteOneByName(ctx context.Context, name values.BreedName) error {
	res, err := b.db.Delete(goqu.T("breeds")).Where(goqu.C("name").Eq(name)).Executor().ExecContext(ctx)
	if err != nil {
		return domainerror.WrapError(domainerror.ErrInternalError, err)
	}
	if n, err := res.RowsAffected(); err != nil {
		return domainerror.WrapError(domainerror.ErrInternalError, err)
	} else if n == 0 {
		return domainerror.WrapError(domainerror.ErrResourceNotFound, fmt.Errorf("breed %s not found", name))
	}
	return nil
}

func (b BreedStorage) List(ctx context.Context, params breeds.ListOpts) ([]*breeds.Breed, error) {
	var res []BreedModel

	query := b.db.From("breeds").
		Select(
			goqu.C("name"),
			goqu.C("pet_size"),
			goqu.C("average_male_adult_weight"),
			goqu.C("average_female_adult_weight"),
			goqu.C("species"),
		)
	if params.Species != nil {
		query = query.Where(goqu.C("species").Eq(params.Species.String()))
	}
	if params.AverageFemaleWeight != nil {
		query = query.Where(goqu.C("average_female_adult_weight").Eq(*params.AverageFemaleWeight))
	}
	if params.AverageMaleWeight != nil {
		query = query.Where(goqu.C("average_male_adult_weight").Eq(*params.AverageMaleWeight))
	}
	if err := query.ScanStructsContext(ctx, &res); err != nil {
		return nil, domainerror.WrapError(domainerror.ErrInternalError, err)
	}
	return common.EMap(res, func(val BreedModel) (*breeds.Breed, error) {
		return val.ToDomain()
	})
}
