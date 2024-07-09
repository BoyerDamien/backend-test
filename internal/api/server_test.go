package api_test

import (
	"context"
	"fmt"
	"net/http"
	"testing"

	charmLog "github.com/charmbracelet/log"
	"github.com/gorilla/mux"
	"github.com/japhy-tech/backend-test/internal/api"
	"github.com/japhy-tech/backend-test/internal/common"
	"github.com/japhy-tech/backend-test/internal/domain/breeds"
	"github.com/japhy-tech/backend-test/internal/domain/values"
	"github.com/japhy-tech/backend-test/internal/domainerror"
	"github.com/japhy-tech/backend-test/internal/gateways"
	"github.com/japhy-tech/backend-test/internal/testutils"
	"github.com/japhy-tech/backend-test/internal/usecases"
	breedUsecases "github.com/japhy-tech/backend-test/internal/usecases/breeds"
	"github.com/maxatome/go-testdeep/helpers/tdhttp"
	"github.com/maxatome/go-testdeep/td"
)

func TestServer_CreateOneBreed(t *testing.T) {
	testutils.TestDecorator(t, func(ctx context.Context, datastore gateways.IDatastore, require *td.T, logger *charmLog.Logger) {
		var (
			r     = mux.NewRouter()
			h     = api.HandlerFromMuxWithBaseURL(api.New(logger, datastore), r, "/v1")
			ta    = tdhttp.NewTestAPI(t, h)
			tests = []struct {
				name           string
				body           api.Breed
				expectedStatus int
				errContains    string
			}{
				{
					name: "valid case",
					body: api.Breed{
						Name:                     "test",
						Species:                  api.Species(values.Cat.String()),
						PetSize:                  api.PetSize(values.Medium.String()),
						AverageFemaleAdultWeight: common.ToPointer(1),
						AverageMaleAdultWeight:   common.ToPointer(1),
					},
					expectedStatus: http.StatusCreated,
				},
				{
					name: "invalid case -- already exist",
					body: api.Breed{
						Name:    "test",
						Species: api.Species(values.Cat.String()),
						PetSize: api.PetSize(values.Medium.String()),
					},
					expectedStatus: http.StatusConflict,
					errContains:    domainerror.ErrResourceAlreadyExists.Error(),
				},
				{
					name: "invalid case -- invalid name",
					body: api.Breed{
						Name:                     "test not valid",
						Species:                  api.Species(values.Cat.String()),
						PetSize:                  api.PetSize(values.Medium.String()),
						AverageFemaleAdultWeight: common.ToPointer(1),
						AverageMaleAdultWeight:   common.ToPointer(1),
					},
					expectedStatus: http.StatusBadRequest,
					errContains:    values.ErrNameInvalid.Error(),
				},
				{
					name: "invalid case -- invalid name -- too long",
					body: api.Breed{
						Name:                     common.GenString(256),
						Species:                  api.Species(values.Cat.String()),
						PetSize:                  api.PetSize(values.Medium.String()),
						AverageFemaleAdultWeight: common.ToPointer(1),
						AverageMaleAdultWeight:   common.ToPointer(1),
					},
					expectedStatus: http.StatusBadRequest,
					errContains:    values.ErrNameToLong.Error(),
				},
				{
					name: "invalid case -- invalid name -- too short",
					body: api.Breed{
						Name:                     "o",
						Species:                  api.Species(values.Cat.String()),
						PetSize:                  api.PetSize(values.Medium.String()),
						AverageFemaleAdultWeight: common.ToPointer(1),
						AverageMaleAdultWeight:   common.ToPointer(1),
					},
					expectedStatus: http.StatusBadRequest,
					errContains:    values.ErrNameToShort.Error(),
				},
				{
					name: "invalid case -- invalid species",
					body: api.Breed{
						Name:                     "test_valid",
						Species:                  "api.Species(values.Cat.String())",
						PetSize:                  api.PetSize(values.Medium.String()),
						AverageFemaleAdultWeight: common.ToPointer(1),
						AverageMaleAdultWeight:   common.ToPointer(1),
					},
					expectedStatus: http.StatusBadRequest,
					errContains:    values.ErrInvalidSpecies.Error(),
				},
				{
					name: "invalid case -- invalid petsize",
					body: api.Breed{
						Name:                     "test_valid",
						Species:                  api.Species(values.Cat.String()),
						PetSize:                  "api.PetSize(values.Medium.String())",
						AverageFemaleAdultWeight: common.ToPointer(1),
						AverageMaleAdultWeight:   common.ToPointer(1),
					},
					expectedStatus: http.StatusBadRequest,
					errContains:    values.ErrInvalidPetSize.Error(),
				},
			}
		)

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				ta.Name(tt.name).PostJSON("/v1/breeds", tt.body).
					CmpStatus(tt.expectedStatus)
				if tt.errContains != "" {
					ta.CmpJSONBody(td.JSON(`{"message": $message}`, td.Tag("message", td.Contains(tt.errContains))))
				} else {
					ta.CmpJSONBody(tt.body)
				}
			})
		}
	})
}

func TestServer_CreateOrUpdateBreedByName(t *testing.T) {
	testutils.TestDecorator(t, func(ctx context.Context, datastore gateways.IDatastore, require *td.T, logger *charmLog.Logger) {
		var (
			r     = mux.NewRouter()
			h     = api.HandlerFromMuxWithBaseURL(api.New(logger, datastore), r, "/v1")
			ta    = tdhttp.NewTestAPI(t, h)
			tests = []struct {
				name           string
				body           api.Breed
				expectedStatus int
				errContains    string
			}{
				{
					name: "valid case",
					body: api.Breed{
						Name:                     "test",
						Species:                  api.Species(values.Cat.String()),
						PetSize:                  api.PetSize(values.Medium.String()),
						AverageFemaleAdultWeight: common.ToPointer(1),
						AverageMaleAdultWeight:   common.ToPointer(1),
					},
					expectedStatus: http.StatusCreated,
				},
				{
					name: "invalid case -- already exist",
					body: api.Breed{
						Name:                     "test",
						Species:                  api.Species(values.Cat.String()),
						PetSize:                  api.PetSize(values.Medium.String()),
						AverageFemaleAdultWeight: common.ToPointer(1),
						AverageMaleAdultWeight:   common.ToPointer(1),
					},
					expectedStatus: http.StatusNoContent,
				},
				{
					name: "valid case -- update species",
					body: api.Breed{
						Name:                     "test",
						Species:                  api.Species(values.Dog.String()),
						PetSize:                  api.PetSize(values.Medium.String()),
						AverageFemaleAdultWeight: common.ToPointer(1),
						AverageMaleAdultWeight:   common.ToPointer(1),
					},
					expectedStatus: http.StatusOK,
				},
				{
					name: "valid case -- update petsize",
					body: api.Breed{
						Name:                     "test",
						Species:                  api.Species(values.Dog.String()),
						PetSize:                  api.PetSize(values.Tall.String()),
						AverageFemaleAdultWeight: common.ToPointer(1),
						AverageMaleAdultWeight:   common.ToPointer(1),
					},
					expectedStatus: http.StatusOK,
				},
				{
					name: "valid case -- update weight",
					body: api.Breed{
						Name:                     "test",
						Species:                  api.Species(values.Dog.String()),
						PetSize:                  api.PetSize(values.Tall.String()),
						AverageFemaleAdultWeight: common.ToPointer(10),
						AverageMaleAdultWeight:   common.ToPointer(30),
					},
					expectedStatus: http.StatusOK,
				},
				{
					name: "invalid case -- invalid name",
					body: api.Breed{
						Name:                     "test_not_valid2",
						Species:                  api.Species(values.Cat.String()),
						PetSize:                  api.PetSize(values.Medium.String()),
						AverageFemaleAdultWeight: common.ToPointer(1),
						AverageMaleAdultWeight:   common.ToPointer(1),
					},
					expectedStatus: http.StatusBadRequest,
					errContains:    values.ErrNameInvalid.Error(),
				},
				{
					name: "invalid case -- invalid name -- too long",
					body: api.Breed{
						Name:                     common.GenString(256),
						Species:                  api.Species(values.Cat.String()),
						PetSize:                  api.PetSize(values.Medium.String()),
						AverageFemaleAdultWeight: common.ToPointer(1),
						AverageMaleAdultWeight:   common.ToPointer(1),
					},
					expectedStatus: http.StatusBadRequest,
					errContains:    values.ErrNameToLong.Error(),
				},
				{
					name: "invalid case -- invalid name -- too short",
					body: api.Breed{
						Name:                     "o",
						Species:                  api.Species(values.Cat.String()),
						PetSize:                  api.PetSize(values.Medium.String()),
						AverageFemaleAdultWeight: common.ToPointer(1),
						AverageMaleAdultWeight:   common.ToPointer(1),
					},
					expectedStatus: http.StatusBadRequest,
					errContains:    values.ErrNameToShort.Error(),
				},
				{
					name: "invalid case -- invalid species",
					body: api.Breed{
						Name:                     "test_valid",
						Species:                  "api.Species(values.Cat.String())",
						PetSize:                  api.PetSize(values.Medium.String()),
						AverageFemaleAdultWeight: common.ToPointer(1),
						AverageMaleAdultWeight:   common.ToPointer(1),
					},
					expectedStatus: http.StatusBadRequest,
					errContains:    values.ErrInvalidSpecies.Error(),
				},
				{
					name: "invalid case -- invalid petsize",
					body: api.Breed{
						Name:                     "test_valid",
						Species:                  api.Species(values.Cat.String()),
						PetSize:                  "api.PetSize(values.Medium.String())",
						AverageFemaleAdultWeight: common.ToPointer(1),
						AverageMaleAdultWeight:   common.ToPointer(1),
					},
					expectedStatus: http.StatusBadRequest,
					errContains:    values.ErrInvalidPetSize.Error(),
				},
			}
		)

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				ta = ta.Name(tt.name).PutJSON(fmt.Sprintf("/v1/breeds/name/%s", tt.body.Name), tt.body).CmpStatus(tt.expectedStatus)
				if tt.expectedStatus == http.StatusCreated || tt.expectedStatus == http.StatusOK {
					ta.CmpJSONBody(tt.body)
				}
				if tt.errContains != "" && tt.expectedStatus != http.StatusNoContent {
					ta.CmpJSONBody(td.JSON(`{"message": $message}`, td.Tag("message", td.Contains(tt.errContains))))
				}
			})
		}
	})
}

func TestServer_GetBreedByName(t *testing.T) {
	testutils.TestDecorator(t, func(ctx context.Context, datastore gateways.IDatastore, require *td.T, logger *charmLog.Logger) {
		var (
			createHandler = usecases.New(&breedUsecases.CreateOne{}, datastore)
			r             = mux.NewRouter()
			h             = api.HandlerFromMuxWithBaseURL(api.New(logger, datastore), r, "/v1")
			ta            = tdhttp.NewTestAPI(t, h)

			tests = []struct {
				name           string
				input          string
				expectedStatus int
				errContains    string
			}{
				{
					name:           "valid case",
					input:          "test",
					expectedStatus: http.StatusOK,
				},
				{
					name:           "invalid case -- not found",
					input:          "not_found",
					expectedStatus: http.StatusNotFound,
				},
				{
					name:           "invalid case -- invalid name",
					input:          "invalid_name2",
					expectedStatus: http.StatusBadRequest,
					errContains:    values.ErrNameInvalid.Error(),
				},
				{
					name:           "invalid case -- invalid name -- too long",
					input:          common.GenString(256),
					expectedStatus: http.StatusBadRequest,
					errContains:    values.ErrNameToLong.Error(),
				},
				{
					name:           "invalid case -- invalid name -- too short",
					input:          "o",
					expectedStatus: http.StatusBadRequest,
					errContains:    values.ErrNameToShort.Error(),
				},
			}
		)

		_, err := createHandler.Handle(ctx, breeds.FactoryOpts{
			Name:                "test",
			Species:             values.Cat.String(),
			PetSize:             values.Medium.String(),
			AverageFemaleWeight: common.ToPointer(1),
			AverageMaleWeight:   common.ToPointer(1),
		})
		require.CmpNoError(err)

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				ta = ta.Name(tt.name).Get(fmt.Sprintf("/v1/breeds/name/%s", tt.input)).CmpStatus(tt.expectedStatus)
				if tt.expectedStatus == http.StatusOK {
					ta.CmpJSONBody(td.JSON(
						`{
						"name": "test",
						"species": "cat",
						"pet_size": "medium",
						"average_female_adult_weight": 1,
						"average_male_adult_weight": 1,
						}`))
				} else {
					ta.CmpJSONBody(td.JSON(`{"message": $message}`, td.Tag("message", td.Contains(tt.errContains))))
				}
			})
		}
	})
}

func TestServer_DeleteBreedByName(t *testing.T) {
	testutils.TestDecorator(t, func(ctx context.Context, datastore gateways.IDatastore, require *td.T, logger *charmLog.Logger) {
		var (
			createHandler = usecases.New(&breedUsecases.CreateOne{}, datastore)
			r             = mux.NewRouter()
			h             = api.HandlerFromMuxWithBaseURL(api.New(logger, datastore), r, "/v1")
			ta            = tdhttp.NewTestAPI(t, h)

			tests = []struct {
				name           string
				input          string
				expectedStatus int
				errContains    string
			}{
				{
					name:           "valid case",
					input:          "test",
					expectedStatus: http.StatusNoContent,
				},
				{
					name:           "invalid case -- not found",
					input:          "not_found",
					expectedStatus: http.StatusNotFound,
				},
				{
					name:           "invalid case -- invalid name",
					input:          "invalid_name2",
					expectedStatus: http.StatusBadRequest,
					errContains:    values.ErrNameInvalid.Error(),
				},
				{
					name:           "invalid case -- invalid name -- too long",
					input:          common.GenString(256),
					expectedStatus: http.StatusBadRequest,
					errContains:    values.ErrNameToLong.Error(),
				},
				{
					name:           "invalid case -- invalid name -- too short",
					input:          "o",
					expectedStatus: http.StatusBadRequest,
					errContains:    values.ErrNameToShort.Error(),
				},
			}
		)

		_, err := createHandler.Handle(ctx, breeds.FactoryOpts{
			Name:                "test",
			Species:             values.Cat.String(),
			PetSize:             values.Medium.String(),
			AverageFemaleWeight: common.ToPointer(1),
			AverageMaleWeight:   common.ToPointer(1),
		})
		require.CmpNoError(err)

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				ta = ta.Name(tt.name).Delete(fmt.Sprintf("/v1/breeds/name/%s", tt.input), nil).CmpStatus(tt.expectedStatus)
				if tt.expectedStatus != http.StatusNoContent {
					ta.CmpJSONBody(td.JSON(`{"message": $message}`, td.Tag("message", td.Contains(tt.errContains))))
				}
			})
		}
	})
}

func TestServer_ListBreeds(t *testing.T) {
	testutils.TestDecorator(t, func(ctx context.Context, datastore gateways.IDatastore, require *td.T, logger *charmLog.Logger) {
		var (
			r         = mux.NewRouter()
			h         = api.HandlerFromMuxWithBaseURL(api.New(logger, datastore), r, "/v1")
			ta        = tdhttp.NewTestAPI(t, h)
			breedArgs = []breeds.FactoryOpts{
				{
					Name:                "test_dog",
					Species:             values.Dog.String(),
					PetSize:             values.Medium.String(),
					AverageFemaleWeight: common.ToPointer(10),
					AverageMaleWeight:   common.ToPointer(3),
				},
				{
					Name:                "test_cat",
					Species:             values.Cat.String(),
					PetSize:             values.Medium.String(),
					AverageFemaleWeight: common.ToPointer(10),
					AverageMaleWeight:   common.ToPointer(1),
				},
				{
					Name:                "test_cat_small",
					Species:             values.Cat.String(),
					PetSize:             values.Small.String(),
					AverageFemaleWeight: common.ToPointer(2),
					AverageMaleWeight:   common.ToPointer(1),
				},
				{
					Name:                "test_dog_tall",
					Species:             values.Dog.String(),
					PetSize:             values.Tall.String(),
					AverageFemaleWeight: common.ToPointer(3),
					AverageMaleWeight:   common.ToPointer(1),
				},
				{
					Name:    "test_dog_no_weight",
					Species: values.Dog.String(),
					PetSize: values.Medium.String(),
				},
			}

			breedsCreated []api.Breed
		)

		for _, val := range breedArgs {
			b, err := breeds.NewFactory(val).Instantiate()
			require.CmpNoError(err)
			b, err = datastore.Breeds().CreateOne(context.Background(), b)
			require.CmpNoError(err)
			breedsCreated = append(breedsCreated, api.BreedToJson(b))
		}

		tests := []struct {
			name           string
			queryFilter    string
			expectedResult []api.Breed
			expectedStatus int
		}{
			{
				name:           "get all",
				expectedResult: breedsCreated,
				expectedStatus: http.StatusOK,
			},
			{
				name: "get dog only",
				expectedResult: []api.Breed{
					breedsCreated[0],
					breedsCreated[3],
					breedsCreated[4],
				},
				queryFilter:    "?species=dog",
				expectedStatus: http.StatusOK,
			},
			{
				name: "get cat only",
				expectedResult: []api.Breed{
					breedsCreated[1],
					breedsCreated[2],
				},
				queryFilter:    "?species=cat",
				expectedStatus: http.StatusOK,
			},
			{
				name: "get medium only",
				expectedResult: []api.Breed{
					breedsCreated[0],
					breedsCreated[1],
					breedsCreated[4],
				},
				queryFilter:    "?pet_size=medium",
				expectedStatus: http.StatusOK,
			},
			{
				name: "get small only",
				expectedResult: []api.Breed{
					breedsCreated[2],
				},
				queryFilter:    "?pet_size=small",
				expectedStatus: http.StatusOK,
			},
			{
				name: "get tall only",
				expectedResult: []api.Breed{
					breedsCreated[3],
				},
				expectedStatus: http.StatusOK,
				queryFilter:    "?pet_size=tall",
			},

			{
				name: "get male weight 0",
				expectedResult: []api.Breed{
					breedsCreated[4],
				},
				expectedStatus: http.StatusOK,
				queryFilter:    "?average_male_adult_weight=0",
			},
			{
				name: "get female weight 0",
				expectedResult: []api.Breed{
					breedsCreated[4],
				},
				expectedStatus: http.StatusOK,
				queryFilter:    "?average_female_adult_weight=0",
			},
			{
				name: "get female weight 2",
				expectedResult: []api.Breed{
					breedsCreated[2],
				},
				expectedStatus: http.StatusOK,
				queryFilter:    "?average_female_adult_weight=2",
			},
			{
				name: "get male weight 1",
				expectedResult: []api.Breed{
					breedsCreated[1],
					breedsCreated[2],
					breedsCreated[3],
				},
				expectedStatus: http.StatusOK,
				queryFilter:    "?average_male_adult_weight=1",
			},
			{
				name: "get dog and medium",
				expectedResult: []api.Breed{
					breedsCreated[0],
					breedsCreated[4],
				},
				expectedStatus: http.StatusOK,
				queryFilter:    "?pet_size=medium&species=dog",
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				ta.Name(tt.name).Get("/v1/breeds" + tt.queryFilter).CmpStatus(tt.expectedStatus).CmpJSONBody(tt.expectedResult)
			})
		}
	})
}
