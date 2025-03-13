package repository

import (
	"fmt"
	"webPractice1/internal/domain"
)

func (c *CRUD) GetEntity(ip string) (*domain.AssetData, error) {
	asset := domain.NewAsset()
	asset.Mu.Lock()
	defer asset.Mu.Unlock()
	tx, err := c.db.Begin()
	if err != nil {
		c.logger.Error(fmt.Sprintf("transaction not started: %s", err))
		return nil, err
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
		c.logger.Error(fmt.Sprintf("SELECT Entity ERROR: %s", err))
		return nil, err
	}
	defer query.Close()
	for query.Next() {
		if err := query.Scan(&asset.Asset.Id, &asset.Asset.IPAddress, &asset.Asset.IsPublic, &asset.Asset.IPVersion, &asset.Asset.IsWhitelisted,
			&asset.Asset.AbuseConfidenceScore, &asset.Asset.CountryCode, &asset.Asset.CountryName, &asset.Asset.UsageType,
			&asset.Asset.ISP, &asset.Asset.IsTor, &asset.Asset.IsDb); err != nil {
			c.logger.Error(fmt.Sprintf("Scan method error: %s", err))
			return nil, err
		}
	}
	if asset.Asset.Id == "" {
		return nil, domain.ErrNoEntityFound
	}
	return &asset.Asset, nil
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
		if err := query.Scan(&asset.Id, &asset.IPAddress, &asset.IsPublic, &asset.IPVersion, &asset.IsWhitelisted,
			&asset.AbuseConfidenceScore, &asset.CountryCode, &asset.CountryName, &asset.UsageType,
			&asset.ISP, &asset.IsTor, &asset.IsDb); err != nil {
			c.logger.Error(fmt.Sprintf("Scan method error: %s", err))
		}
		assets = append(assets, asset)
	}
	return assets
}

func (c *CRUD) GetEntityById(ip string) (int, error) {
	tx, err := c.db.Begin()
	if err != nil {
		c.logger.Error(fmt.Sprintf("Truntransaction not started: %s", err))
		return 0, err
	}
	defer func() {
		if err != nil {
			tx.Rollback()
			c.logger.Error(fmt.Sprintf("Something wrong with transaction: %s", err))
		} else {
			tx.Commit()
		}
	}()
	query := tx.QueryRow(`SELECT "Id" FROM "`+c.crudDb+`" WHERE "ipAddress"=$1`, ip)
	var id int
	if err := query.Scan(&id); err != nil {
		c.logger.Error(fmt.Sprintf("Scan ID ERROR: %s", err))
		return 0, err
	}
	return id, err
}
