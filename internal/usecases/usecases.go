package usecases

import (
	"context"
	"fmt"
	"strings"

	"github.com/japhy-tech/backend-test/internal/gateways"
	"github.com/japhy-tech/backend-test/internal/logger"
)

type UsecaseAction int
type UsecaseName int

const (
	ActionCreate UsecaseAction = iota
	ActionDelete
	ActionUpdate
	ActionRetrieve
	ActionList

	BreedUsecase UsecaseName = iota
)

func (u UsecaseAction) String() string {
	switch u {
	case ActionDelete:
		return "delete"
	case ActionRetrieve:
		return "retrieve"
	case ActionUpdate:
		return "update"
	case ActionCreate:
		return "create"
	case ActionList:
		return "list"
	default:
		return ""
	}
}

func (u UsecaseName) String() string {
	switch u {
	case BreedUsecase:
		return "breed"
	default:
		return ""
	}
}

type UseCaseInfo struct {
	Action UsecaseAction
	Name   UsecaseName
}

func (u UseCaseInfo) String() string {
	return fmt.Sprintf("<%s %s>", strings.ToUpper(u.Action.String()), strings.ToUpper(u.Name.String()))
}

type IBase interface {
	Init(gateways.IDatastore)
	Datastore() gateways.IDatastore
	Info() UseCaseInfo
}

type IUsecase[Input any, Output any] interface {
	Handle(context.Context, Input) (Output, error)
	IBase
}

type ISimpleUsecase[Input any] interface {
	Handle(context.Context, Input) error
	IBase
}

type Default[Input any, Output any] struct {
	content IUsecase[Input, Output]
}

func (b Default[Input, Output]) Handle(ctx context.Context, input Input) (Output, error) {
	l := logger.Logger
	l.Infof("Execute usecase %s", b.content.Info())

	r, err := b.content.Handle(ctx, input)
	if err != nil {
		l.Errorf("Usecase %s [FAILED]", b.content.Info())
	} else {
		l.Infof("Usecase %s [SUCCEED]", b.content.Info())
	}
	return r, err
}

func (b Default[Input, Output]) Init(datastore gateways.IDatastore) {
	b.content.Init(datastore)
}

func (b Default[Input, Output]) Datastore() gateways.IDatastore {
	return b.content.Datastore()
}

func (b Default[Input, Output]) Info() UseCaseInfo {
	return b.content.Info()
}

type SimpleDefault[Input any] struct {
	content ISimpleUsecase[Input]
}

func (b SimpleDefault[Input]) Handle(ctx context.Context, input Input) error {
	l := logger.Logger
	l.Infof("Execute usecase %s", b.content.Info())

	err := b.content.Handle(ctx, input)
	if err != nil {
		l.Errorf("Usecase %s [FAILED]", b.content.Info())
	} else {
		l.Infof("Usecase %s [SUCCEED]", b.content.Info())
	}
	return err
}
func (b SimpleDefault[Input]) Init(datastore gateways.IDatastore) {
	b.content.Init(datastore)
}

func (b SimpleDefault[Input]) Datastore() gateways.IDatastore {
	return b.content.Datastore()
}

func (b SimpleDefault[Input]) Info() UseCaseInfo {
	return b.content.Info()
}

func New[Input any, Output any](usecase IUsecase[Input, Output], datastore gateways.IDatastore) IUsecase[Input, Output] {
	r := &Default[Input, Output]{content: usecase}
	r.Init(datastore)
	return r
}

func NewSimple[Input any](usecase ISimpleUsecase[Input], datastore gateways.IDatastore) ISimpleUsecase[Input] {
	r := &SimpleDefault[Input]{content: usecase}
	r.Init(datastore)
	return r
}
