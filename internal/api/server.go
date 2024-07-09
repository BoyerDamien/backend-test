package api

import (
	"net/http"

	charmLog "github.com/charmbracelet/log"
)

// Server
// Implement ServerInterface
type Server struct {
	logger *charmLog.Logger
}

// List breeds
// (GET /breeds)
func (s Server) ListBreeds(w http.ResponseWriter, r *http.Request, params ListBreedsParams) {
	panic("not implemented") // TODO: Implement
}

// Create one breed
// (POST /breeds)
func (s Server) CreateOneBreed(w http.ResponseWriter, r *http.Request) {
	panic("not implemented") // TODO: Implement
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

func New(logger *charmLog.Logger) *Server {
	return &Server{
		logger: logger,
	}
}
