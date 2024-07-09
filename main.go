//go:generate go run github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen --config=config.yaml ./api/swagger.yml

package main

import (
	"net"
	"net/http"
	"time"

	charmLog "github.com/charmbracelet/log"
	"github.com/gorilla/mux"
	"github.com/japhy-tech/backend-test/internal/api"
	"github.com/japhy-tech/backend-test/internal/gateways/mysql"
	"github.com/japhy-tech/backend-test/internal/logger"
)

const (
	MysqlDSN = "root:root@(mysql-test:3306)/core?parseTime=true"
	ApiPort  = "5000"
)

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

func main() {

	datastore := mysql.New(MysqlDSN, logger.Logger)
	defer datastore.Close()

	r := mux.NewRouter()
	r.Use(loggingMiddleware(logger.Logger))
	r.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}).Methods(http.MethodGet)

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
