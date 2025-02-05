package repository

import (
	"fmt"
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

func (c *CRUD) DeleteEntityDB(ip string) {
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
	_, err = tx.Exec(`DELETE FROM "`+c.crudDb+`" WHERE "ipAddress" = $1`, ip)
	if err != nil {
		c.logger.Error(fmt.Sprintf("DELETE IN DB ERROR: %s", err))
		return
	}
	if err != nil {
		return
	}
}
