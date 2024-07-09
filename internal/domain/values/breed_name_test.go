package values_test

import (
	"testing"

	"github.com/japhy-tech/backend-test/internal/domain/values"
	"github.com/maxatome/go-testdeep/td"
)

func TestBreedName_Validate(t *testing.T) {
	tests := []struct {
		name    string
		b       values.BreedName
		wantErr error
	}{
		{
			name: "valid case -- with underscore",
			b:    "test_ok",
		},
		{
			name: "valid case -- without underscore",
			b:    "test",
		},
		{
			name:    "invalid case -- with space",
			b:       "test invalid",
			wantErr: values.ErrNameInvalid,
		},
	}

	require := td.Require(t)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.b.Validate()
			require.CmpErrorIs(err, tt.wantErr)
		})
	}
}
