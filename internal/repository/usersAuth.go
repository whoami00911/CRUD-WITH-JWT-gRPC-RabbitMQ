package repository

import (
	"fmt"
	"webPractice1/internal/domain"
)

func (ad *AuthDatabase) CreateUser(user domain.User) (int, error) {
	tx, err := ad.db.Begin()
	if err != nil {
		ad.logger.Error(fmt.Sprintf("transaction not started: %s", err))
		return 0, err
	}
	defer func() {
		if err != nil {
			tx.Rollback()
			ad.logger.Error(fmt.Sprintf("Something wrong with transaction: %s", err))
		} else {
			tx.Commit()
		}
	}()
	var id int
	row := tx.QueryRow(`INSERT INTO "`+ad.usersDb+`" ("name", "username", "password_hash") VALUES ($1, $2, $3) RETURNING "id"`, user.Name, user.Username, user.Password)
	if err := row.Scan(&id); err != nil {
		ad.logger.Error(fmt.Sprintf("Select Scan method error: %s", err))
		return 0, err
	}
	return id, err
}

func (ad *AuthDatabase) GetUser(user, password string) (int, error) {
	tx, err := ad.db.Begin()
	if err != nil {
		ad.logger.Error(fmt.Sprintf("transaction not started: %s", err))
		return 0, err
	}
	defer func() {
		if err != nil {
			tx.Rollback()
			ad.logger.Error(fmt.Sprintf("Something wrong with transaction: %s", err))
		} else {
			tx.Commit()
		}
	}()
	var id int
	row := tx.QueryRow(`SELECT "id" FROM "`+ad.usersDb+`" WHERE "username" = $1 AND "password_hash" = $2`, user, password)
	if err := row.Scan(&id); err != nil {
		ad.logger.Error(fmt.Sprintf("Get Scan method error: %s", err))
		return 0, err
	}
	return id, nil
}
