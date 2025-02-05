package repository

import (
	"database/sql"
	"fmt"
	"webPractice1/internal/domain"
	"webPractice1/pkg/logger"

	"github.com/spf13/viper"
)

type SessionDb struct {
	db          *sql.DB
	logger      *logger.Logger
	tableTokens string
}

func NewSessionDb(db *sql.DB, log *logger.Logger) *SessionDb {
	return &SessionDb{
		db:          db,
		logger:      log,
		tableTokens: viper.GetString("db_tables.tokens"),
	}
}
func (sd *SessionDb) CreateRToken(token domain.RefreshSession) {
	tx, err := sd.db.Begin()
	if err != nil {
		sd.logger.Error(fmt.Sprintf("transaction not started: %s", err))
		return
	}
	defer func() {
		if err != nil {
			tx.Rollback()
			sd.logger.Error(fmt.Sprintf("Something wrong with transaction: %s", err))
		} else {
			tx.Commit()
		}
	}()
	if _, err := tx.Exec(`INSERT INTO "`+sd.tableTokens+`" ("userId", "refreshToken", "expiresAt") VALUES ($1, $2, $3) RETURNING "id"`, token.UserID, token.Token, token.ExpiresAt); err != nil {
		sd.logger.Error(fmt.Sprintf("INSERT TOKEN ERROR: %s", err))
	}

}
func (sd *SessionDb) GetRToken(token string) (domain.RefreshSession, error) {
	tx, err := sd.db.Begin()
	if err != nil {
		sd.logger.Error(fmt.Sprintf("transaction not started: %s", err))
		return domain.RefreshSession{}, err
	}
	defer func() {
		if err != nil {
			tx.Rollback()
			sd.logger.Error(fmt.Sprintf("Something wrong with transaction: %s", err))
		} else {
			tx.Commit()
		}
	}()
	var parameters domain.RefreshSession
	row := tx.QueryRow(`SELECT * FROM "`+sd.tableTokens+`" WHERE "refreshToken" = $1`, token)
	if err := row.Scan(&parameters.ID, &parameters.UserID, &parameters.Token, &parameters.ExpiresAt); err != nil {
		sd.logger.Error(fmt.Sprintf("Get Scan method error: %s", err))
		return domain.RefreshSession{}, err
	}
	tx.Exec(`DELETE FROM "`+sd.tableTokens+`" WHERE "userId" = $1`, parameters.UserID)
	//sd.logger.Info(fmt.Sprintf("Token expiry in DB: %v", parameters.ExpiresAt))
	return parameters, nil
}
