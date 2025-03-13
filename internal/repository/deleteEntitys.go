package repository

import (
	"fmt"
	"webPractice1/internal/domain"
)

func (c *CRUD) DeleteAllEntitiesDB() {
	tx, err := c.db.Begin()
	if err != nil {
		c.logger.Error(fmt.Sprintf("transaction not started: %s", err))
		return
	}
	defer func() {
		if err != nil {
			tx.Rollback()
			c.logger.Error(fmt.Sprintf("Something wrong with transaction: %s", err))
			return
		} else {
			tx.Commit()
		}
	}()
	_, err = tx.Exec(`DELETE FROM "` + c.crudDb + `"`)
	if err != nil {
		c.logger.Error(fmt.Sprintf("DELETE IN DB ERROR: %s", err))
		return
	}
}

func (c *CRUD) DeleteEntityDB(ip string) error {
	tx, err := c.db.Begin()
	if err != nil {
		c.logger.Error(fmt.Sprintf("transaction not started: %s", err))
		return err
	}
	defer func() {
		if err != nil {
			tx.Rollback()
			c.logger.Error(fmt.Sprintf("Something wrong with transaction: %s", err))
			return
		} else {
			tx.Commit()
		}
	}()
	result, err := tx.Exec(`DELETE FROM "`+c.crudDb+`" WHERE "ipAddress" = $1`, ip)
	if err != nil {
		c.logger.Error(fmt.Sprintf("DELETE IN DB ERROR: %s", err))
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		c.logger.Error(fmt.Sprintf("Error fetching rows affected: %s", err))
		return err
	}

	if rowsAffected == 0 {
		return domain.ErrNoEntityFound
	}

	return nil
}
