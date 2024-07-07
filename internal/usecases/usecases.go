package usecases

import (
	"context"
	"fmt"

	"github.com/japhy-tech/backend-test/internal/gateways"
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
	return fmt.Sprintf("%s - %s", u.Action, u.Name)
}

type IUsecase[Input any, Output any] interface {
	Handle(context.Context, Input) (Output, error)
	Init(gateways.IDatastore)
	Datastore() gateways.IDatastore
	Info()
}

func New[Input any, Output any](usecase IUsecase[Input, Output], datastore gateways.IDatastore) IUsecase[Input, Output] {
	usecase.Init(datastore)
	return usecase
}
