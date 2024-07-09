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

func TestGetOneByName_Handle(t *testing.T) {
	testutils.TestDecorator(t, func(ctx context.Context, datastore gateways.IDatastore, require *td.T, _ *charmLog.Logger) {
		var (
			createHandler = usecases.New(&breedUsecases.CreateOne{}, datastore)
			getHandler    = usecases.New(&breedUsecases.GetOneByName{}, datastore)

			tests = []struct {
				name        string
				input       string
				wantErr     error
				errContains string
			}{
				{
					name:  "valid case",
					input: "test",
				},
				{
					name:    "invalid case -- not found",
					input:   "not_found",
					wantErr: domainerror.ErrResourceNotFound,
				},
				{
					name:        "invalid case -- invalid name",
					input:       "invalid name",
					wantErr:     domainerror.ErrDomainValidation,
					errContains: values.ErrNameInvalid.Error(),
				},
				{
					name:        "invalid case -- invalid name -- too long",
					input:       common.GenString(256),
					wantErr:     domainerror.ErrDomainValidation,
					errContains: values.ErrNameToLong.Error(),
				},
				{
					name:        "invalid case -- invalid name -- too short",
					input:       "o",
					wantErr:     domainerror.ErrDomainValidation,
					errContains: values.ErrNameToShort.Error(),
				},
			}
		)

		b, err := createHandler.Handle(ctx, breeds.FactoryOpts{
			Name:                "test",
			Species:             values.Cat.String(),
			PetSize:             values.Medium.String(),
			AverageFemaleWeight: common.ToPointer(1),
			AverageMaleWeight:   common.ToPointer(1),
		})
		require.CmpNoError(err)

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				res, err := getHandler.Handle(ctx, tt.input)
				require.CmpErrorIs(err, tt.wantErr)

				if tt.wantErr != nil {
					require.Contains(err.Error(), tt.errContains)
				} else {
					require.Cmp(res, b)
				}
			})
		}
	})
}
