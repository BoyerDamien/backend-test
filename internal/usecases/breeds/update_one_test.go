package breeds_test

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
	"github.com/japhy-tech/backend-test/internal/usecases"
	breedUsecases "github.com/japhy-tech/backend-test/internal/usecases/breeds"
	"github.com/maxatome/go-testdeep/td"
)

func TestUpdateOne_Handle(t *testing.T) {
	testutils.TestDecorator(t, func(ctx context.Context, datastore gateways.IDatastore, require *td.T, _ *charmLog.Logger) {
		var (
			createHandler = usecases.New(&breedUsecases.CreateOne{}, datastore)
			updateHandler = usecases.New(&breedUsecases.UpdateOne{}, datastore)

			tests = []struct {
				name        string
				input       breeds.FactoryOpts
				wantErr     error
				errContains string
			}{
				{
					name: "no change",
					input: breeds.FactoryOpts{
						Name:                "test",
						Species:             values.Cat.String(),
						PetSize:             values.Medium.String(),
						AverageFemaleWeight: common.ToPointer(1),
						AverageMaleWeight:   common.ToPointer(1),
					},
					wantErr: domainerror.ErrNothingTodo,
				},
				{
					name: "valid case -- without weight",
					input: breeds.FactoryOpts{
						Name:    "test",
						Species: values.Cat.String(),
						PetSize: values.Medium.String(),
					},
				},
				{
					name: "invalid case -- invalid name",
					input: breeds.FactoryOpts{
						Name:                "test not valid",
						Species:             values.Cat.String(),
						PetSize:             values.Medium.String(),
						AverageFemaleWeight: common.ToPointer(1),
						AverageMaleWeight:   common.ToPointer(1),
					},
					wantErr:     domainerror.ErrDomainValidation,
					errContains: values.ErrNameInvalid.Error(),
				},
				{
					name: "invalid case -- invalid name -- too long",
					input: breeds.FactoryOpts{
						Name:                common.GenString(256),
						Species:             values.Cat.String(),
						PetSize:             values.Medium.String(),
						AverageFemaleWeight: common.ToPointer(1),
						AverageMaleWeight:   common.ToPointer(1),
					},
					wantErr:     domainerror.ErrDomainValidation,
					errContains: values.ErrNameToLong.Error(),
				},
				{
					name: "invalid case -- invalid name -- too short",
					input: breeds.FactoryOpts{
						Name:                "o",
						Species:             values.Cat.String(),
						PetSize:             values.Medium.String(),
						AverageFemaleWeight: common.ToPointer(1),
						AverageMaleWeight:   common.ToPointer(1),
					},
					wantErr:     domainerror.ErrDomainValidation,
					errContains: values.ErrNameToShort.Error(),
				},
				{
					name: "invalid case -- invalid species",
					input: breeds.FactoryOpts{
						Name:                "test",
						Species:             "values.Cat.String()",
						PetSize:             values.Medium.String(),
						AverageFemaleWeight: common.ToPointer(1),
						AverageMaleWeight:   common.ToPointer(1),
					},
					wantErr:     domainerror.ErrDomainValidation,
					errContains: values.ErrInvalidSpecies.Error(),
				},
				{
					name: "invalid case -- invalid petsize",
					input: breeds.FactoryOpts{
						Name:                "test",
						Species:             values.Cat.String(),
						PetSize:             "values.Medium.String()",
						AverageFemaleWeight: common.ToPointer(1),
						AverageMaleWeight:   common.ToPointer(1),
					},
					wantErr:     domainerror.ErrDomainValidation,
					errContains: values.ErrInvalidPetSize.Error(),
				},
			}
		)

		_, err := createHandler.Handle(ctx, breeds.FactoryOpts{
			Name:                "test",
			Species:             values.Cat.String(),
			PetSize:             values.Medium.String(),
			AverageFemaleWeight: common.ToPointer(1),
			AverageMaleWeight:   common.ToPointer(1),
		})
		require.CmpNoError(err)

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				res, err := updateHandler.Handle(ctx, tt.input)
				require.CmpErrorIs(err, tt.wantErr)

				if tt.wantErr != nil {
					require.Contains(err.Error(), tt.errContains)
				} else {
					require.Cmp(res.Name().String(), tt.input.Name)
					require.Cmp(res.PetSize().String(), tt.input.PetSize)
					require.Cmp(res.Species().String(), tt.input.Species)

					if tt.input.AverageFemaleWeight != nil {
						require.Cmp(res.AverageFemaleWeight(), *tt.input.AverageFemaleWeight)
					} else {
						require.Cmp(res.AverageFemaleWeight(), 0)
					}
					if tt.input.AverageMaleWeight != nil {
						require.Cmp(res.AverageMaleWeight(), *tt.input.AverageMaleWeight)
					} else {
						require.Cmp(res.AverageMaleWeight(), 0)
					}
				}
			})
		}
	})
}
