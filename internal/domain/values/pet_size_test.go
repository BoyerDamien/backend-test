package values_test

import (
	"strings"
	"testing"

	"github.com/japhy-tech/backend-test/internal/domain/values"
	"github.com/maxatome/go-testdeep/td"
)

func TestPetSizeFromString(t *testing.T) {
	tests := []struct {
		name    string
		s       string
		want    values.PetSize
		wantErr error
	}{
		{
			name: "valid case -- small",
			s:    values.Small.String(),
			want: values.Small,
		},
		{
			name: "valid case -- medium",
			s:    values.Medium.String(),
			want: values.Medium,
		},
		{
			name: "valid case -- tall",
			s:    values.Tall.String(),
			want: values.Tall,
		},
		{
			name: "valid case -- small uppercase",
			s:    strings.ToUpper(values.Small.String()),
			want: values.Small,
		},
		{
			name: "valid case -- medium uppercase",
			s:    strings.ToUpper(values.Medium.String()),
			want: values.Medium,
		},
		{
			name: "valid case -- tall uppercase",
			s:    strings.ToUpper(values.Tall.String()),
			want: values.Tall,
		},
		{
			name:    "invalid case",
			s:       "invalid",
			wantErr: values.ErrInvalidPetSize,
		},
		{
			name:    "invalid case -- empty string",
			s:       "",
			wantErr: values.ErrInvalidPetSize,
		},
	}

	require := td.Require(t)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := values.PetSizeFromString(tt.s)
			require.CmpErrorIs(err, tt.wantErr)
			if tt.wantErr == nil {
				require.Cmp(got, tt.want)
			}
		})
	}
}
