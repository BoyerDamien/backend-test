package values

import (
	"errors"
	"fmt"
	"regexp"
)

type BreedName string

const (
	BreedNameRegexp = `^[a-z]+(_[a-z]+)*$`
)

var (
	ErrNameToShort = errors.New("breed name must have at least 2 characters")
	ErrNameToLong  = errors.New("breed name must have maximum 255 characters")
	ErrNameInvalid = fmt.Errorf("breed name does not follow this pattern %s", BreedNameRegexp)
)

// Validate
// Implements values.Validator interface
func (b BreedName) Validate() error {
	l := len(b)

	if l < 2 {
		return ErrNameToShort
	}
	if l > 255 {
		return ErrNameToLong
	}

	// Regexp compile should not failed here
	match, _ := regexp.MatchString(BreedNameRegexp, b.String())
	if !match {
		return ErrNameInvalid
	}
	return nil
}

func (b BreedName) String() string {
	return string(b)
}
