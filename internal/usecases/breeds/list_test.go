package breeds_test

import (
	"context"
	"testing"

	"github.com/japhy-tech/backend-test/internal/common"
	"github.com/japhy-tech/backend-test/internal/domain/breeds"
	"github.com/japhy-tech/backend-test/internal/domain/values"
	"github.com/japhy-tech/backend-test/internal/gateways"
	"github.com/japhy-tech/backend-test/internal/testutils"
	"github.com/japhy-tech/backend-test/internal/usecases"
	breedUsecases "github.com/japhy-tech/backend-test/internal/usecases/breeds"
	"github.com/maxatome/go-testdeep/td"
)

func TestList_Handle(t *testing.T) {
	testutils.TestDecorator(t, func(ctx context.Context, datastore gateways.IDatastore, require *td.T) {
		var (
			handler   = usecases.New(&breedUsecases.List{}, datastore)
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
			filter         breedUsecases.ListOpts
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
				filter: breedUsecases.ListOpts{
					Species: common.ToPointer(values.Dog.String()),
				},
			},
			{
				name: "get cat only",
				expectedResult: []*breeds.Breed{
					breedsCreated[1],
					breedsCreated[2],
				},
				filter: breedUsecases.ListOpts{
					Species: common.ToPointer(values.Cat.String()),
				},
			},
			{
				name: "get medium only",
				expectedResult: []*breeds.Breed{
					breedsCreated[0],
					breedsCreated[1],
					breedsCreated[4],
				},
				filter: breedUsecases.ListOpts{
					PetSize: common.ToPointer(values.Medium.String()),
				},
			},
			{
				name: "get small only",
				expectedResult: []*breeds.Breed{
					breedsCreated[2],
				},
				filter: breedUsecases.ListOpts{
					PetSize: common.ToPointer(values.Small.String()),
				},
			},
			{
				name: "get tall only",
				expectedResult: []*breeds.Breed{
					breedsCreated[3],
				},
				filter: breedUsecases.ListOpts{
					PetSize: common.ToPointer(values.Tall.String()),
				},
			},

			{
				name: "get male weight 0",
				expectedResult: []*breeds.Breed{
					breedsCreated[4],
				},
				filter: breedUsecases.ListOpts{
					AverageMaleWeight: common.ToPointer(0),
				},
			},
			{
				name: "get female weight 0",
				expectedResult: []*breeds.Breed{
					breedsCreated[4],
				},
				filter: breedUsecases.ListOpts{
					AverageFemaleWeight: common.ToPointer(0),
				},
			},
			{
				name: "get female weight 2",
				expectedResult: []*breeds.Breed{
					breedsCreated[2],
				},
				filter: breedUsecases.ListOpts{
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
				filter: breedUsecases.ListOpts{
					AverageMaleWeight: common.ToPointer(1),
				},
			},
			{
				name: "get dog and medium",
				expectedResult: []*breeds.Breed{
					breedsCreated[0],
					breedsCreated[4],
				},
				filter: breedUsecases.ListOpts{
					PetSize: common.ToPointer(values.Medium.String()),
					Species: common.ToPointer(values.Dog.String()),
				},
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				res, err := handler.Handle(ctx, tt.filter)
				require.CmpNoError(err)
				require.Cmp(res, tt.expectedResult)
			})
		}
	})
}
