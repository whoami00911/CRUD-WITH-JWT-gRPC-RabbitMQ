package service

import (
	"webPractice1/internal/domain"
	"webPractice1/internal/repository"
	"webPractice1/pkg/hasher"
	"webPractice1/pkg/logger"
)

type Autherization interface {
	CreateUser(user domain.User) (int, error)
}

type Session interface {
	GenTokens(user, password string) (string, string, error)
	ParseToken(token string) (int, error)
	CreateRToken(token domain.RefreshSession)
	GetRToken(token string) (domain.RefreshSession, error)
	UpdateTokens(token string) (string, string, error)
}

type CRUDList interface {
	AddEntity(ar domain.AssetData)
	DeleteAllEntitiesDB()
	DeleteEntityDB(ip string)
	GetEntity(ip string) *domain.AssetData
	GetEntities() []domain.AssetData
	UpdateEntity(ar domain.AssetData)
}

type Service struct {
	Autherization
	CRUDList
	Session
}

func NewService(repos *repository.Repository, hash *hasher.Hash, logger *logger.Logger) *Service {
	return &Service{
		Autherization: NewAuthService(repos.Authorization, hash),
		CRUDList:      NewServiceCRUD(repos.CRUDList),
		Session:       newSessionRepo(repos.Session, repos.Authorization, hash, logger),
	}
}
