package values_test

import (
	"strings"
	"testing"

	"github.com/japhy-tech/backend-test/internal/domain/values"
	"github.com/maxatome/go-testdeep/td"
)

func TestSpeciesFromString(t *testing.T) {
	tests := []struct {
		name    string
		s       string
		want    values.Species
		wantErr error
	}{
		{
			name: "valid case -- cat",
			s:    values.Cat.String(),
			want: values.Cat,
		},
		{
			name: "valid case -- dog",
			s:    values.Dog.String(),
			want: values.Dog,
		},
		{
			name: "valid case -- cat uppercase",
			s:    strings.ToUpper(values.Cat.String()),
			want: values.Cat,
		},
		{
			name: "valid case -- dog uppercase",
			s:    strings.ToUpper(values.Dog.String()),
			want: values.Dog,
		},
		{
			name:    "invalid case",
			s:       "invalid",
			wantErr: values.ErrInvalidSpecies,
		},
		{
			name:    "invalid case -- empty string",
			s:       "",
			wantErr: values.ErrInvalidSpecies,
		},
	}

	require := td.Require(t)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := values.SpeciesFromString(tt.s)
			require.CmpErrorIs(err, tt.wantErr)
			if tt.wantErr == nil {
				require.Cmp(got, tt.want)
			}
		})
	}
}
