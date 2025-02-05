package repository

import (
	"fmt"
	"webPractice1/internal/domain"
)

func (c *CRUD) AddEntity(ar domain.AssetData) {
	tx, err := c.db.Begin()
	if err != nil {
		c.logger.Error(fmt.Sprintf("Truntransaction not started: %s", err))
		return
	}
	defer func() {
		if err != nil {
			tx.Rollback()
			c.logger.Error(fmt.Sprintf("Something wrong with truntransaction: %s", err))
		} else {
			tx.Commit()
		}
	}()
	ar.IsDb = true
	_, err = tx.Exec(
		`INSERT INTO "`+c.crudDb+`" ("ipAddress", "isPublic", "ipVersion", "isWhitelisted", "abuseConfidenceScore", "countryCode", "countryName", "usageType", "isFromDB", "isTor", "isp")
		 VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11) ON CONFLICT ("ipAddress")
		DO NOTHING`, ar.IPAddress, ar.IsPublic, ar.IPVersion, ar.IsWhitelisted, ar.AbuseConfidenceScore,
		ar.CountryCode, ar.CountryName, ar.UsageType, ar.IsDb, ar.IsTor, ar.ISP)
	if err != nil {
		c.logger.Error(fmt.Sprintf("INSERT ERROR: %s", err))
		return
	}
}
