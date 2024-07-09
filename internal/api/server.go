package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

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

// List breeds
// (GET /breeds)
func (s Server) ListBreeds(w http.ResponseWriter, r *http.Request, params ListBreedsParams) {
	panic("not implemented") // TODO: Implement
}

func HandleErrorResponse(w http.ResponseWriter, err error) {
	var errMap = map[string]int{
		domainerror.ErrDomainValidation.Error():      http.StatusBadRequest,
		domainerror.ErrInternalError.Error():         http.StatusInternalServerError,
		domainerror.ErrResourceAlreadyExists.Error(): http.StatusConflict,
		domainerror.ErrNothingTodo.Error():           http.StatusNoContent,
		domainerror.ErrResourceNotFound.Error():      http.StatusNotFound,
	}

	if err.Error() == "EOF" {
		_ = SendJSON(w, Error{Message: domainerror.WrapError(domainerror.ErrDomainValidation, errors.New("body is required")).Error()}, http.StatusBadRequest)
		return
	}

	unwrapped := errors.Unwrap(err)
	val, ok := errMap[err.Error()]

	if unwrapped == nil && !ok {
		_ = SendJSON(w, Error{Message: domainerror.ErrInternalError.Error()}, http.StatusInternalServerError)
		return
	}

	if unwrapped == nil && ok {
		_ = SendJSON(w, Error{Message: err.Error()}, val)
		return
	}
	_ = SendJSON(w, Error{Message: err.Error()}, errMap[unwrapped.Error()])
}

// Create one breed
// (POST /breeds)
func (s Server) CreateOneBreed(w http.ResponseWriter, r *http.Request) {
	var body Breed

	err := json.NewDecoder(r.Body).Decode(&body)
	fmt.Println("err = ", err, "body = ", body)
	if err != nil {
		HandleErrorResponse(w, err)
		return
	}

	res, err := usecases.New(&breedsUsecase.CreateOne{}, s.datastore).Handle(r.Context(), breeds.FactoryOpts{
		Name:                body.Name,
		Species:             string(body.Species),
		PetSize:             string(body.PetSize),
		AverageFemaleWeight: body.AverageFemaleAdultWeight,
		AverageMaleWeight:   body.AverageMaleAdultWeight,
	})
	if err != nil {
		HandleErrorResponse(w, err)
		return
	}

	SendJSON(w, Breed{
		Name:                     res.Name().String(),
		Species:                  Species(res.Species().String()),
		PetSize:                  BreedsPetSize(res.PetSize().String()),
		AverageFemaleAdultWeight: common.ToPointer(res.AverageFemaleWeight()),
		AverageMaleAdultWeight:   common.ToPointer(res.AverageMaleWeight()),
	}, http.StatusCreated)
}

// Delete a given breed by its name
// (DELETE /breeds/name/{breed_name})
func (s Server) DeleteBreedByName(w http.ResponseWriter, r *http.Request, breedName BreedName) {
	panic("not implemented") // TODO: Implement
}

// Retrieve a given breed by its name
// (GET /breeds/name/{breed_name})
func (s Server) GetBreedByName(w http.ResponseWriter, r *http.Request, breedName BreedName) {
	panic("not implemented") // TODO: Implement
}

// Update or create one breed
// (PUT /breeds/name/{breed_name})
func (s Server) CreateOrUpdateBreedByName(w http.ResponseWriter, r *http.Request, breedName BreedName) {
	panic("not implemented") // TODO: Implement
}

func New(logger *charmLog.Logger, datastore gateways.IDatastore) *Server {
	return &Server{
		logger:    logger,
		datastore: datastore,
	}
}
