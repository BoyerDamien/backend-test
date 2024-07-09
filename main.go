//go:generate go run github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen --config=config.yaml ./api/swagger.yml

package main

import (
	"context"
	"encoding/csv"
	"errors"
	"fmt"
	"net"
	"net/http"
	"os"
	"strconv"
	"time"

	charmLog "github.com/charmbracelet/log"
	"github.com/gorilla/mux"
	"github.com/japhy-tech/backend-test/internal/api"
	"github.com/japhy-tech/backend-test/internal/common"
	"github.com/japhy-tech/backend-test/internal/domain/breeds"
	"github.com/japhy-tech/backend-test/internal/domain/values"
	"github.com/japhy-tech/backend-test/internal/domainerror"
	"github.com/japhy-tech/backend-test/internal/gateways"
	"github.com/japhy-tech/backend-test/internal/gateways/mysql"
	"github.com/japhy-tech/backend-test/internal/logger"
)

const (
	MysqlDSN = "root:root@(mysql-test:3306)/core?parseTime=true"
	ApiPort  = "5000"
)

func main() {
	// Init datastore
	datastore := mysql.New(MysqlDSN, logger.Logger)
	defer datastore.Close()

	/// Sync data from csv with the datastore
	breeds, err := breedsFromCSV("./breeds.csv")
	if err != nil {
		logger.Logger.Fatalf("cannot convert csv data: %s", err)
	}
	if err := syncDatastore(breeds, datastore); err != nil {
		logger.Logger.Fatalf("cannot insert csv data in datastore: %s", err)
	}

	// Init Api handler
	r := mux.NewRouter()
	r.Use(loggingMiddleware(logger.Logger))
	r.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}).Methods(http.MethodGet)
	r.PathPrefix("/v1/docs/").Handler(http.StripPrefix("/v1/docs/", http.FileServer(http.Dir("./api"))))

	h := api.HandlerFromMuxWithBaseURL(api.New(logger.Logger, datastore), r, "/v1")

	server := &http.Server{
		Handler: h,
		Addr:    net.JoinHostPort("", ApiPort),
	}

	if err := server.ListenAndServe(); err != nil {
		logger.Logger.Fatal(err.Error())
	}

	// =============================== Starting Msg ===============================
	logger.Logger.Infof("Service started and listen on port %s", ApiPort)
}

func loggingMiddleware(logger *charmLog.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			next.ServeHTTP(w, r)
			logger.Infof(
				"%s | %s | %s | %s | %s",
				r.Method,
				r.URL.Path,
				time.Now().Format(time.RFC822),
				r.RemoteAddr,
				r.UserAgent(),
			)
		})
	}
}

func readCsvFile(filePath string) [][]string {
	logger.Logger.Infof("starting reading %s", filePath)
	f, err := os.Open(filePath)
	if err != nil {
		logger.Logger.Fatal("Unable to read input file "+filePath, err)
	}
	defer f.Close()

	csvReader := csv.NewReader(f)
	records, err := csvReader.ReadAll()
	if err != nil {
		logger.Logger.Fatal("Unable to parse file as CSV for "+filePath, err)
	}
	return records
}

func breedsFromCSV(filepath string) ([]*breeds.Breed, error) {
	var (
		records = readCsvFile(filepath)
		res     []*breeds.Breed
	)
	if len(records) > 1 {
		records = records[1:]
	}

	for _, val := range records {
		averageMale, err := strconv.Atoi(val[4])
		if err != nil {
			return nil, fmt.Errorf("cannot convert average male adult weight value %s from csv file %s with records %+v", val[4], filepath, val)
		}
		averageFemale, err := strconv.Atoi(val[5])
		if err != nil {
			return nil, fmt.Errorf("cannot convert average female adult weight value %s from csv file %s with records %+v", val[5], filepath, val)
		}
		b, err := breeds.NewFactory(breeds.FactoryOpts{
			Name:                val[3],
			PetSize:             val[2],
			Species:             val[1],
			AverageFemaleWeight: &averageFemale,
			AverageMaleWeight:   &averageMale,
		}).Instantiate()
		if err != nil {
			return nil, err
		}
		res = append(res, b)
	}
	logger.Logger.Infof("%d elements were found", len(res))
	return res, nil
}

func syncDatastore(arr []*breeds.Breed, datastore gateways.IDatastore) error {
	var (
		ctx          = context.Background()
		breedNameArr = common.Map(arr, func(val *breeds.Breed) string {
			return val.Name().String()
		})
		memo     = make(map[values.BreedName]*breeds.Breed)
		toInsert []*breeds.Breed
		ErrFail  = errors.New("fail to synchronize datastore")
	)

	logger.Logger.Info("Stating datastore synchronization")
	found, err := datastore.Breeds().List(ctx, breeds.ListOpts{NameIn: breedNameArr})
	if err != nil {
		return domainerror.WrapError(ErrFail, err)
	}

	for _, val := range found {
		memo[val.Name()] = val
	}
	for _, val := range arr {
		_, ok := memo[val.Name()]
		if !ok {
			toInsert = append(toInsert, val)
		}
	}

	logger.Logger.Infof("%d elements will be inserted", len(toInsert))
	if len(toInsert) > 0 {
		_, err = datastore.Breeds().CreateSeveral(ctx, toInsert)
		if err != nil {
			return domainerror.WrapError(ErrFail, err)
		}
	}
	logger.Logger.Info("datastore synchronized")
	return nil
}
