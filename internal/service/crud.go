package service

import (
	"webPractice1/internal/domain"
	"webPractice1/internal/repository"
)

type repositoryCRUD struct {
	repo repository.CRUDList
}

func NewServiceCRUD(repo repository.CRUDList) *repositoryCRUD {
	return &repositoryCRUD{
		repo: repo,
	}
}

func (rc *repositoryCRUD) AddEntity(ar domain.AssetData) error {
	return rc.repo.AddEntity(ar)
}
func (rc *repositoryCRUD) DeleteAllEntitiesDB() {
	rc.repo.DeleteAllEntitiesDB()
}
func (rc *repositoryCRUD) DeleteEntityDB(ip string) error {
	return rc.repo.DeleteEntityDB(ip)
}
func (rc *repositoryCRUD) GetEntity(ip string) (*domain.AssetData, error) {
	return rc.repo.GetEntity(ip)
}
func (rc *repositoryCRUD) GetEntities() []domain.AssetData {
	return rc.repo.GetEntities()
}
func (rc *repositoryCRUD) UpdateEntity(ar domain.AssetData) error {
	return rc.repo.UpdateEntity(ar)
}
func (rc *repositoryCRUD) GetEntityById(ip string) (int, error) {
	return rc.repo.GetEntityById(ip)
}
