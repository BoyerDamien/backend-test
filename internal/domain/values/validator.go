package values

import (
	"github.com/japhy-tech/backend-test/internal/domainerror"
)

func Verify(validators ...Validator) error {
	var errArr []error
	for _, v := range validators {
		if err := v.Validate(); err != nil {
			errArr = append(errArr, err)
		}
	}
	if len(errArr) > 0 {
		return domainerror.WrapError(domainerror.ErrDomainValidation, errArr...)
	}
	return nil
}
