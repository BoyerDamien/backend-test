package api

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	charmLog "github.com/charmbracelet/log"
	"github.com/japhy-tech/backend-test/internal/common"
	"github.com/japhy-tech/backend-test/internal/domain/breeds"
	"github.com/japhy-tech/backend-test/internal/domainerror"
	"github.com/japhy-tech/backend-test/internal/gateways"
	"github.com/japhy-tech/backend-test/internal/usecases"
	breedsUsecase "github.com/japhy-tech/backend-test/internal/usecases/breeds"
)

// Server
// Implement ServerInterface
type Server struct {
	logger    *charmLog.Logger
	datastore gateways.IDatastore
}

type Response[T any] struct {
	Status int
	Val    T
}

// List breeds
// (GET /breeds)
func (s Server) ListBreeds(w http.ResponseWriter, r *http.Request, params ListBreedsParams) {
	EndpointDecorator(w, r, func(ctx context.Context) (*Response[[]Breed], error) {
		res, err := usecases.New(&breedsUsecase.List{}, s.datastore).Handle(ctx, breedsUsecase.ListOpts{
			Species:             (*string)(params.Species),
			AverageFemaleWeight: params.AverageFemaleAdultWeight,
			AverageMaleWeight:   params.AverageMaleAdultWeight,
		})
		if err != nil {
			return nil, err
		}
		return &Response[[]Breed]{
			Val:    common.Map(res, func(val *breeds.Breed) Breed { return BreedToJson(val) }),
			Status: http.StatusOK,
		}, nil
	})
}

// Create one breed
// (POST /breeds)
func (s Server) CreateOneBreed(w http.ResponseWriter, r *http.Request) {
	EndpointDecorator(w, r, func(ctx context.Context) (*Response[Breed], error) {
		body, err := Bind[Breed](r)
		if err != nil {
			return nil, err
		}

		res, err := usecases.New(&breedsUsecase.CreateOne{}, s.datastore).Handle(ctx, breeds.FactoryOpts{
			Name:                body.Name,
			Species:             string(body.Species),
			PetSize:             string(body.PetSize),
			AverageFemaleWeight: body.AverageFemaleAdultWeight,
			AverageMaleWeight:   body.AverageMaleAdultWeight,
		})
		if err != nil {
			return nil, err
		}
		return &Response[Breeds]{
			Val:    BreedToJson(res),
			Status: http.StatusCreated,
		}, nil
	})
}

// Delete a given breed by its name
// (DELETE /breeds/name/{breed_name})
func (s Server) DeleteBreedByName(w http.ResponseWriter, r *http.Request, breedName BreedName) {
	err := usecases.New(&breedsUsecase.DeleteOneByName{}, s.datastore).Handle(r.Context(), breedName)
	if err != nil {
		HandleErrorResponse(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// Retrieve a given breed by its name
// (GET /breeds/name/{breed_name})
func (s Server) GetBreedByName(w http.ResponseWriter, r *http.Request, breedName BreedName) {
	EndpointDecorator(w, r, func(ctx context.Context) (*Response[Breed], error) {
		res, err := usecases.New(&breedsUsecase.GetOneByName{}, s.datastore).Handle(ctx, breedName)
		if err != nil {
			return nil, err
		}

		return &Response[Breeds]{
			Val:    BreedToJson(res),
			Status: http.StatusOK,
		}, nil
	})
}

// Update or create one breed
// (PUT /breeds/name/{breed_name})
func (s Server) CreateOrUpdateBreedByName(w http.ResponseWriter, r *http.Request, breedName BreedName) {
	EndpointDecorator(w, r, func(ctx context.Context) (*Response[Breed], error) {
		body, err := Bind[Breed](r)
		if err != nil {
			return nil, err
		}

		opts := breeds.FactoryOpts{
			Name:                breedName,
			Species:             string(body.Species),
			PetSize:             string(body.PetSize),
			AverageFemaleWeight: body.AverageFemaleAdultWeight,
			AverageMaleWeight:   body.AverageMaleAdultWeight,
		}

		res, err := usecases.New(&breedsUsecase.CreateOne{}, s.datastore).Handle(ctx, opts)
		if err == nil {
			return &Response[Breeds]{
				Val:    BreedToJson(res),
				Status: http.StatusCreated,
			}, nil
		}
		res, err = usecases.New(&breedsUsecase.UpdateOne{}, s.datastore).Handle(ctx, opts)
		if err != nil {
			return nil, err
		}

		return &Response[Breeds]{
			Val:    BreedToJson(res),
			Status: http.StatusOK,
		}, nil
	})
}

func New(logger *charmLog.Logger, datastore gateways.IDatastore) *Server {
	return &Server{
		logger:    logger,
		datastore: datastore,
	}
}

func Bind[T any](r *http.Request) (T, error) {
	var body T
	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		return body, domainerror.WrapError(domainerror.ErrDomainValidation, err)
	}
	return body, nil
}

func HandleErrorResponse(w http.ResponseWriter, err error) {
	var errMap = map[string]int{
		domainerror.ErrDomainValidation.Error():      http.StatusBadRequest,
		domainerror.ErrInternalError.Error():         http.StatusInternalServerError,
		domainerror.ErrResourceAlreadyExists.Error(): http.StatusConflict,
		domainerror.ErrNothingTodo.Error():           http.StatusNoContent,
		domainerror.ErrResourceNotFound.Error():      http.StatusNotFound,
	}

	if strings.Contains(err.Error(), "EOF") {
		_ = SendJSON(w, Error{Message: domainerror.WrapError(domainerror.ErrDomainValidation, errors.New("body is required")).Error()}, http.StatusBadRequest)
		return
	}

	unwrapped := errors.Unwrap(err)
	val, ok := errMap[err.Error()]

	if unwrapped == nil && !ok {
		_ = SendJSON(w, Error{Message: err.Error()}, http.StatusInternalServerError)
		return
	}

	if unwrapped == nil && ok {
		_ = SendJSON(w, Error{Message: err.Error()}, val)
		return
	}
	_ = SendJSON(w, Error{Message: err.Error()}, errMap[unwrapped.Error()])
}

func BreedToJson(domain *breeds.Breed) Breed {
	return Breeds{
		Name:                     domain.Name().String(),
		Species:                  Species(domain.Species().String()),
		PetSize:                  BreedsPetSize(domain.PetSize().String()),
		AverageFemaleAdultWeight: common.ToPointer(domain.AverageFemaleWeight()),
		AverageMaleAdultWeight:   common.ToPointer(domain.AverageMaleWeight()),
	}
}

func EndpointDecorator[Output any](w http.ResponseWriter, r *http.Request, fn func(context.Context) (*Response[Output], error)) {
	res, err := fn(r.Context())
	if err != nil {
		HandleErrorResponse(w, err)
	} else {
		SendJSON(w, res.Val, res.Status)
	}
}
