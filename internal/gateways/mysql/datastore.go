package mysql

import (
	"database/sql"
	"os"

	charmLog "github.com/charmbracelet/log"
	"github.com/doug-martin/goqu/v9"
	_ "github.com/doug-martin/goqu/v9/dialect/mysql"
	"github.com/japhy-tech/backend-test/database_actions"
	"github.com/japhy-tech/backend-test/internal/domain/breeds"
)

type Datastore struct {
	breeds *BreedStorage
	logger *charmLog.Logger
	db     *sql.DB
}

func (d Datastore) Close() error {
	return d.db.Close()
}

func (d Datastore) Breeds() breeds.Repository {
	return d.breeds
}

func New(dsn string, logger *charmLog.Logger) *Datastore {
	err := database_actions.InitMigrator(dsn)
	if err != nil {
		logger.Fatal(err.Error())
	}

	msg, err := database_actions.RunMigrate("up", 0)
	if err != nil {
		logger.Error(err.Error())
	} else {
		logger.Info(msg)
	}

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		logger.Fatal(err.Error())
		os.Exit(1)
	}
	db.SetMaxIdleConns(0)

	err = db.Ping()
	if err != nil {
		logger.Fatal(err.Error())
		os.Exit(1)
	}

	logger.Info("Database connected")

	return &Datastore{
		breeds: NewBreedStorage(goqu.New("mysql", db)),
		db:     db,
		logger: logger,
	}
}
