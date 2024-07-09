package mysql_test

import (
	"context"
	"testing"

	charmLog "github.com/charmbracelet/log"
	"github.com/japhy-tech/backend-test/internal/common"
	"github.com/japhy-tech/backend-test/internal/domain/breeds"
	"github.com/japhy-tech/backend-test/internal/domain/values"
	"github.com/japhy-tech/backend-test/internal/domainerror"
	"github.com/japhy-tech/backend-test/internal/gateways"
	"github.com/japhy-tech/backend-test/internal/testutils"
	"github.com/maxatome/go-testdeep/td"
)

func TestBreedStorage_CreateOne(t *testing.T) {
	testutils.TestDecorator(t, func(ctx context.Context, datastore gateways.IDatastore, require *td.T, _ *charmLog.Logger) {
		type Args struct {
			name                     values.BreedName
			species                  values.Species
			petSize                  values.PetSize
			averageFemaleAdultWeight *int
			averageMaleAdultWeight   *int
		}

		var (
			tests = []struct {
				name    string
				args    Args
				wantErr error
			}{
				{
					name: "valid case",
					args: Args{
						name:                     "test",
						species:                  values.Cat,
						petSize:                  values.Medium,
						averageFemaleAdultWeight: common.ToPointer(1),
						averageMaleAdultWeight:   common.ToPointer(1),
					},
				},
				{
					name: "valid case -- without weight",
					args: Args{
						name:    "test_ok",
						species: values.Cat,
						petSize: values.Medium,
					},
				},
				{
					name: "invalid case -- already exist",
					args: Args{
						name:    "test_ok",
						species: values.Cat,
						petSize: values.Medium,
					},
					wantErr: domainerror.ErrInternalError,
				},
			}
		)

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				b, err := breeds.NewFactory(breeds.FactoryOpts{
					Name:                tt.args.name.String(),
					Species:             tt.args.species.String(),
					PetSize:             tt.args.petSize.String(),
					AverageFemaleWeight: tt.args.averageFemaleAdultWeight,
					AverageMaleWeight:   tt.args.averageMaleAdultWeight,
				}).Instantiate()
				require.CmpNoError(err)

				res, err := datastore.Breeds().CreateOne(context.Background(), b)
				require.CmpErrorIs(err, tt.wantErr)

				if tt.wantErr == nil {
					require.Cmp(res, td.Struct(&breeds.Breed{}, td.StructFields{
						"name":    tt.args.name,
						"species": tt.args.species,
						"petSize": tt.args.petSize,
						"averageFemaleWeight": func() int {
							if tt.args.averageFemaleAdultWeight != nil {
								return *tt.args.averageFemaleAdultWeight
							}
							return 0
						}(),
						"averageMaleWeight": func() int {
							if tt.args.averageMaleAdultWeight != nil {
								return *tt.args.averageMaleAdultWeight
							}
							return 0
						}(),
					}))
				}
			})
		}

	})
}

func TestBreedStorage_GetOneByName(t *testing.T) {
	testutils.TestDecorator(t, func(ctx context.Context, datastore gateways.IDatastore, require *td.T, _ *charmLog.Logger) {
		var (
			tests = []struct {
				name      string
				breedName values.BreedName
				wantErr   error
			}{
				{
					name:      "valid case",
					breedName: "test",
				},
				{
					name:      "invalid case -- name not found",
					breedName: "not_found",
					wantErr:   domainerror.ErrResourceNotFound,
				},
			}
		)

		b, err := breeds.NewFactory(breeds.FactoryOpts{
			Name:    "test",
			Species: values.Cat.String(),
			PetSize: values.Small.String(),
		}).Instantiate()
		require.CmpNoError(err)

		res, err := datastore.Breeds().CreateOne(context.Background(), b)
		require.CmpNoError(err)

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				r, err := datastore.Breeds().GetOneByName(ctx, tt.breedName)
				require.CmpErrorIs(err, tt.wantErr)

				if tt.wantErr == nil {
					require.Cmp(r, td.Struct(&breeds.Breed{}, td.StructFields{
						"name":                res.Name(),
						"species":             res.Species(),
						"petSize":             res.PetSize(),
						"averageFemaleWeight": res.AverageFemaleWeight(),
						"averageMaleWeight":   res.AverageMaleWeight(),
					}))
				}
			})
		}
	})
}

func TestBreedStorage_DeleteOneByName(t *testing.T) {
	testutils.TestDecorator(t, func(ctx context.Context, datastore gateways.IDatastore, require *td.T, _ *charmLog.Logger) {
		var (
			tests = []struct {
				name      string
				breedName values.BreedName
				wantErr   error
			}{
				{
					name:      "valid case",
					breedName: "test",
				},
				{
					name:      "invalid case -- name not found",
					breedName: "not_found",
					wantErr:   domainerror.ErrResourceNotFound,
				},
			}
		)

		b, err := breeds.NewFactory(breeds.FactoryOpts{
			Name:    "test",
			Species: values.Cat.String(),
			PetSize: values.Small.String(),
		}).Instantiate()
		require.CmpNoError(err)

		res, err := datastore.Breeds().CreateOne(context.Background(), b)
		require.CmpNoError(err)

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				err := datastore.Breeds().DeleteOneByName(ctx, tt.breedName)
				require.CmpErrorIs(err, tt.wantErr)

				if tt.wantErr == nil {
					_, err := datastore.Breeds().GetOneByName(ctx, res.Name())
					require.CmpErrorIs(err, domainerror.ErrResourceNotFound)
				}
			})
		}
	})
}

func TestBreedStorage_UpdateOne(t *testing.T) {
	testutils.TestDecorator(t, func(ctx context.Context, datastore gateways.IDatastore, require *td.T, _ *charmLog.Logger) {
		b, err := breeds.NewFactory(breeds.FactoryOpts{
			Name:    "test",
			Species: values.Cat.String(),
			PetSize: values.Small.String(),
		}).Instantiate()
		require.CmpNoError(err)

		bUpdated, err := breeds.NewFactory(breeds.FactoryOpts{
			Name:                "test",
			Species:             values.Dog.String(),
			PetSize:             values.Medium.String(),
			AverageFemaleWeight: common.ToPointer(10),
			AverageMaleWeight:   common.ToPointer(1),
		}).Instantiate()
		require.CmpNoError(err)

		_, err = datastore.Breeds().CreateOne(context.Background(), b)
		require.CmpNoError(err)

		r, err := datastore.Breeds().UpdateOne(ctx, bUpdated)
		require.CmpNoError(err)
		require.Cmp(r.Name(), bUpdated.Name())
		require.Cmp(r.Species(), bUpdated.Species())
		require.Cmp(r.AverageFemaleWeight(), bUpdated.AverageFemaleWeight())
		require.Cmp(r.AverageMaleWeight(), bUpdated.AverageMaleWeight())
		require.Cmp(r.PetSize(), bUpdated.PetSize())
	})
}

func TestBreedStorage_List(t *testing.T) {
	testutils.TestDecorator(t, func(ctx context.Context, datastore gateways.IDatastore, require *td.T, _ *charmLog.Logger) {
		var (
			breedArgs = []breeds.FactoryOpts{
				{
					Name:                "test_dog",
					Species:             values.Dog.String(),
					PetSize:             values.Medium.String(),
					AverageFemaleWeight: common.ToPointer(10),
					AverageMaleWeight:   common.ToPointer(3),
				},
				{
					Name:                "test_cat",
					Species:             values.Cat.String(),
					PetSize:             values.Medium.String(),
					AverageFemaleWeight: common.ToPointer(10),
					AverageMaleWeight:   common.ToPointer(1),
				},
				{
					Name:                "test_cat_small",
					Species:             values.Cat.String(),
					PetSize:             values.Small.String(),
					AverageFemaleWeight: common.ToPointer(2),
					AverageMaleWeight:   common.ToPointer(1),
				},
				{
					Name:                "test_dog_tall",
					Species:             values.Dog.String(),
					PetSize:             values.Tall.String(),
					AverageFemaleWeight: common.ToPointer(3),
					AverageMaleWeight:   common.ToPointer(1),
				},
				{
					Name:    "test_dog_no_weight",
					Species: values.Dog.String(),
					PetSize: values.Medium.String(),
				},
			}

			breedsCreated []*breeds.Breed
		)

		for _, val := range breedArgs {
			b, err := breeds.NewFactory(val).Instantiate()
			require.CmpNoError(err)
			b, err = datastore.Breeds().CreateOne(context.Background(), b)
			require.CmpNoError(err)
			breedsCreated = append(breedsCreated, b)
		}

		tests := []struct {
			name           string
			filter         breeds.ListOpts
			expectedResult []*breeds.Breed
		}{
			{
				name:           "get all",
				expectedResult: breedsCreated,
			},
			{
				name: "get dog only",
				expectedResult: []*breeds.Breed{
					breedsCreated[0],
					breedsCreated[3],
					breedsCreated[4],
				},
				filter: breeds.ListOpts{
					Species: common.ToPointer(values.Dog),
				},
			},
			{
				name: "get cat only",
				expectedResult: []*breeds.Breed{
					breedsCreated[1],
					breedsCreated[2],
				},
				filter: breeds.ListOpts{
					Species: common.ToPointer(values.Cat),
				},
			},
			{
				name: "get medium only",
				expectedResult: []*breeds.Breed{
					breedsCreated[0],
					breedsCreated[1],
					breedsCreated[4],
				},
				filter: breeds.ListOpts{
					PetSize: common.ToPointer(values.Medium),
				},
			},
			{
				name: "get small only",
				expectedResult: []*breeds.Breed{
					breedsCreated[2],
				},
				filter: breeds.ListOpts{
					PetSize: common.ToPointer(values.Small),
				},
			},
			{
				name: "get tall only",
				expectedResult: []*breeds.Breed{
					breedsCreated[3],
				},
				filter: breeds.ListOpts{
					PetSize: common.ToPointer(values.Tall),
				},
			},

			{
				name: "get male weight 0",
				expectedResult: []*breeds.Breed{
					breedsCreated[4],
				},
				filter: breeds.ListOpts{
					AverageMaleWeight: common.ToPointer(0),
				},
			},
			{
				name: "get female weight 0",
				expectedResult: []*breeds.Breed{
					breedsCreated[4],
				},
				filter: breeds.ListOpts{
					AverageFemaleWeight: common.ToPointer(0),
				},
			},
			{
				name: "get female weight 2",
				expectedResult: []*breeds.Breed{
					breedsCreated[2],
				},
				filter: breeds.ListOpts{
					AverageFemaleWeight: common.ToPointer(2),
				},
			},
			{
				name: "get male weight 1",
				expectedResult: []*breeds.Breed{
					breedsCreated[1],
					breedsCreated[2],
					breedsCreated[3],
				},
				filter: breeds.ListOpts{
					AverageMaleWeight: common.ToPointer(1),
				},
			},
			{
				name: "get dog and medium",
				expectedResult: []*breeds.Breed{
					breedsCreated[0],
					breedsCreated[4],
				},
				filter: breeds.ListOpts{
					PetSize: common.ToPointer(values.Medium),
					Species: common.ToPointer(values.Dog),
				},
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				res, err := datastore.Breeds().List(ctx, tt.filter)
				require.CmpNoError(err)
				require.Cmp(res, tt.expectedResult)
			})
		}
	})
}
