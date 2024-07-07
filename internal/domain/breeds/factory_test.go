package breeds_test

import (
	"testing"

	"github.com/japhy-tech/backend-test/internal/common"
	"github.com/japhy-tech/backend-test/internal/domain/breeds"
	"github.com/japhy-tech/backend-test/internal/domain/values"
	"github.com/japhy-tech/backend-test/internal/domainerror"
	"github.com/maxatome/go-testdeep/td"
)

func TestFactory_Instantiate(t *testing.T) {
	tests := []struct {
		name        string
		opts        breeds.FactoryOpts
		wantErr     error
		errContains string
	}{
		{
			name: "valid case -- all fields",
			opts: breeds.FactoryOpts{
				Name:                "test",
				Species:             values.Cat.String(),
				PetSize:             values.Small.String(),
				AverageFemaleWeight: common.ToPointer(1000),
				AverageMaleWeight:   common.ToPointer(2000),
			},
		},
		{
			name: "valid case -- without male weight",
			opts: breeds.FactoryOpts{
				Name:              "test",
				Species:           values.Cat.String(),
				PetSize:           values.Small.String(),
				AverageMaleWeight: common.ToPointer(2000),
			},
		},
		{
			name: "valid case -- without femal weight",
			opts: breeds.FactoryOpts{
				Name:              "test",
				Species:           values.Cat.String(),
				PetSize:           values.Small.String(),
				AverageMaleWeight: common.ToPointer(2000),
			},
		},
		{
			name: "invalid case -- invalid name",
			opts: breeds.FactoryOpts{
				Name:                "test invalid",
				Species:             values.Cat.String(),
				PetSize:             values.Small.String(),
				AverageFemaleWeight: common.ToPointer(1000),
				AverageMaleWeight:   common.ToPointer(2000),
			},
			wantErr:     domainerror.ErrDomainValidation,
			errContains: values.ErrNameInvalid.Error(),
		},
		{
			name: "invalid case -- name too short",
			opts: breeds.FactoryOpts{
				Name:                "t",
				Species:             values.Cat.String(),
				PetSize:             values.Small.String(),
				AverageFemaleWeight: common.ToPointer(1000),
				AverageMaleWeight:   common.ToPointer(2000),
			},
			wantErr:     domainerror.ErrDomainValidation,
			errContains: values.ErrNameToShort.Error(),
		},
		{
			name: "invalid case -- name too long",
			opts: breeds.FactoryOpts{
				Name:                common.GenString(256),
				Species:             values.Cat.String(),
				PetSize:             values.Small.String(),
				AverageFemaleWeight: common.ToPointer(1000),
				AverageMaleWeight:   common.ToPointer(2000),
			},
			wantErr:     domainerror.ErrDomainValidation,
			errContains: values.ErrNameToLong.Error(),
		},
		{
			name: "invalid case -- wrong species",
			opts: breeds.FactoryOpts{
				Name:                "test",
				Species:             "invalid",
				PetSize:             values.Small.String(),
				AverageFemaleWeight: common.ToPointer(1000),
				AverageMaleWeight:   common.ToPointer(2000),
			},
			wantErr:     domainerror.ErrDomainValidation,
			errContains: values.ErrInvalidSpecies.Error(),
		},
		{
			name: "valid case -- wrong petsize",
			opts: breeds.FactoryOpts{
				Name:                "test",
				Species:             values.Cat.String(),
				PetSize:             "invalid",
				AverageFemaleWeight: common.ToPointer(1000),
				AverageMaleWeight:   common.ToPointer(2000),
			},
			errContains: values.ErrInvalidPetSize.Error(),
			wantErr:     domainerror.ErrDomainValidation,
		},
	}

	require := td.Require(t)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := breeds.NewFactory(tt.opts)

			got, err := f.Instantiate()
			if tt.wantErr != nil {
				require.CmpErrorIs(err, tt.wantErr)
				require.Contains(err, tt.errContains)
			} else {
				require.Cmp(got.ID(), 0)
				require.Cmp(got.Name().String(), tt.opts.Name)
				require.Cmp(got.PetSize().String(), tt.opts.PetSize)
				if tt.opts.AverageFemaleWeight != nil {
					require.Cmp(got.AverageFemaleWeight(), *tt.opts.AverageFemaleWeight)
				} else {
					require.Cmp(got.AverageFemaleWeight(), 0)
				}
				if tt.opts.AverageMaleWeight != nil {
					require.Cmp(got.AverageMaleWeight(), *tt.opts.AverageMaleWeight)
				} else {
					require.Cmp(got.AverageMaleWeight(), 0)
				}
			}
		})
	}
}
