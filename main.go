//go:generate go run github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen --config=config.yaml ./api/swagger.yml

package main

import (
	"fmt"
	"net"
	"net/http"
	"os"
	"time"

	charmLog "github.com/charmbracelet/log"
	"github.com/gorilla/mux"
	"github.com/japhy-tech/backend-test/internal/api"
	"github.com/japhy-tech/backend-test/internal/gateways/mysql"
)

const (
	MysqlDSN = "root:root@(mysql-test:3306)/core?parseTime=true"
	ApiPort  = "5000"
)

func main() {
	logger := charmLog.NewWithOptions(os.Stderr, charmLog.Options{
		Formatter:       charmLog.TextFormatter,
		ReportCaller:    true,
		ReportTimestamp: true,
		TimeFormat:      time.Kitchen,
		Prefix:          "🧑‍💻 backend-test",
		Level:           charmLog.DebugLevel,
	})

	datastore := mysql.New(MysqlDSN, logger)
	defer datastore.Close()

	r := mux.NewRouter()
	r.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}).Methods(http.MethodGet)

	h := api.HandlerFromMuxWithBaseURL(api.New(logger, datastore), r, "/v1")

	server := &http.Server{
		Handler: h,
		Addr:    net.JoinHostPort("", ApiPort),
	}

	if err := server.ListenAndServe(); err != nil {
		logger.Fatal(err.Error())
	}

	// =============================== Starting Msg ===============================
	logger.Info(fmt.Sprintf("Service started and listen on port %s", ApiPort))
}
