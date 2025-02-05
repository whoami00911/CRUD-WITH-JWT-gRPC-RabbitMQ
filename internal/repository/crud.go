package repository

import (
	"database/sql"
	"webPractice1/pkg/logger"

	"github.com/spf13/viper"
)

type CRUD struct {
	db     *sql.DB
	logger *logger.Logger
	crudDb string
}

func NewCrudDbInicialize(db *sql.DB, log *logger.Logger) *CRUD {
	return &CRUD{
		db:     db,
		logger: log,
		crudDb: viper.GetString("db_tables.crud"),
	}

}
