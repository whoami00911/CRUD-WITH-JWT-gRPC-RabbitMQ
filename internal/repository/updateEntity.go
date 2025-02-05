package repository

import (
	"fmt"
	"webPractice1/internal/domain"
)

func (c *CRUD) UpdateEntity(ar domain.AssetData) {
	tx, err := c.db.Begin()
	if err != nil {
		c.logger.Error(fmt.Sprintf("Truntransaction not started: %s", err))
		return
	}
	defer func() {
		if err != nil {
			tx.Rollback()
			c.logger.Error(fmt.Sprintf("Something wrong with truntransaction: %s", err))
			return
		} else {
			tx.Commit()
		}
	}()
	ar.IsDb = true
	_, err = tx.Exec(
		`UPDATE "`+c.crudDb+`" SET "isPublic" = $1, "ipVersion" = $2, "isWhitelisted" = $3, "abuseConfidenceScore" = $4, "countryCode" = $5, "countryName" = $6, "usageType" = $7, "isFromDB" = $8, "isTor" = $9, "isp" = $10 WHERE "ipAddress" = $11`, ar.IsPublic, ar.IPVersion, ar.IsWhitelisted, ar.AbuseConfidenceScore,
		ar.CountryCode, ar.CountryName, ar.UsageType, ar.IsDb, ar.IsTor, ar.ISP, ar.IPAddress)
	if err != nil {
		c.logger.Error(fmt.Sprintf("UPDATE ERROR: %s", err))
		return
	}

}
