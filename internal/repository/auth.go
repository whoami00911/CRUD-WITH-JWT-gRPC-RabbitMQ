package repository

import (
	"database/sql"
	"webPractice1/pkg/logger"

	"github.com/spf13/viper"
)

type AuthDatabase struct {
	db      *sql.DB
	logger  *logger.Logger
	usersDb string
}

func NewAuthUserDbInicialize(db *sql.DB, log *logger.Logger) *AuthDatabase {
	return &AuthDatabase{
		db:      db,
		logger:  log,
		usersDb: viper.GetString("db_tables.auth"),
	}
}
