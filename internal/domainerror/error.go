package domainerror

import (
	"errors"
	"fmt"
)

var (
	ErrResourceNotFound      = errors.New("resource not found error")
	ErrInternalError         = errors.New("internal error")
	ErrResourceAlreadyExists = errors.New("resource already exists error")
	ErrDomainValidation      = errors.New("resource validation error")
	ErrNothingTodo           = errors.New("nothing to do error")
)

func WrapError(wrapper error, errArr ...error) error {
	return fmt.Errorf("%w: %s", wrapper, errors.Join(errArr...).Error())
}
