package repository

import (
	"fmt"
	"webPractice1/internal/domain"
)

func (c *CRUD) GetEntity(ip string) *domain.AssetData {
	asset := domain.NewAsset()
	asset.Mu.Lock()
	defer asset.Mu.Unlock()
	tx, err := c.db.Begin()
	if err != nil {
		c.logger.Error(fmt.Sprintf("transaction not started: %s", err))
	}
	defer func() {
		if err != nil {
			tx.Rollback() // Rollback only if there was an error
			c.logger.Error(fmt.Sprintf("Something wrong with transaction: %s", err))
			return
		} else {
			tx.Commit() // Commit if no errors occurred
		}
	}()
	query, err := tx.Query(`SELECT * FROM "`+c.crudDb+`" WHERE "ipAddress"=$1`, ip)
	if err != nil {
		c.logger.Error(fmt.Sprintf("INSERT ERROR: %s", err))
	}
	defer query.Close()
	for query.Next() {
		if err := query.Scan(&asset.Asset.IPAddress, &asset.Asset.IsPublic, &asset.Asset.IPVersion, &asset.Asset.IsWhitelisted,
			&asset.Asset.AbuseConfidenceScore, &asset.Asset.CountryCode, &asset.Asset.CountryName, &asset.Asset.UsageType,
			&asset.Asset.ISP, &asset.Asset.IsTor, &asset.Asset.IsDb); err != nil {
			c.logger.Error(fmt.Sprintf("Scan method error: %s", err))
		}
	}
	return &asset.Asset
}
func (c *CRUD) GetEntities() []domain.AssetData {
	var assets []domain.AssetData
	tx, err := c.db.Begin()
	if err != nil {
		c.logger.Error(fmt.Sprintf("Truntransaction not started: %s", err))
	}
	defer func() {
		if err != nil {
			tx.Rollback()
			c.logger.Error(fmt.Sprintf("Something wrong with transaction: %s", err))
		} else {
			tx.Commit()
		}
	}()
	query, err := tx.Query(`SELECT * FROM "` + c.crudDb + `"`)
	if err != nil {
		c.logger.Error(fmt.Sprintf("INSERT ERROR: %s", err))
	}
	defer query.Close()

	for query.Next() {
		var asset domain.AssetData
		if err := query.Scan(&asset.IPAddress, &asset.IsPublic, &asset.IPVersion, &asset.IsWhitelisted,
			&asset.AbuseConfidenceScore, &asset.CountryCode, &asset.CountryName, &asset.UsageType,
			&asset.ISP, &asset.IsTor, &asset.IsDb); err != nil {
			c.logger.Error(fmt.Sprintf("Scan method error: %s", err))
		}
		assets = append(assets, asset)
	}
	return assets
}
