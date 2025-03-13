package service

import (
	"webPractice1/internal/domain"
	"webPractice1/internal/repository"
	"webPractice1/pkg/hasher"
	"webPractice1/pkg/logger"

	"github.com/whoami00911/gRPC-server/pkg/grpcPb"
)

type Autherization interface {
	CreateUser(user domain.User) (int, error)
	GetUserId(user domain.UserSignIn) (int, error)
}

type Session interface {
	GenTokens(user, password string) (string, string, error)
	ParseToken(token string) (int, error)
	CreateRToken(token domain.RefreshSession)
	GetRToken(token string) (domain.RefreshSession, error)
	UpdateTokens(token string) (string, string, error)
}

type CRUDList interface {
	AddEntity(ar domain.AssetData) error
	DeleteAllEntitiesDB()
	DeleteEntityDB(ip string) error
	GetEntity(ip string) (*domain.AssetData, error)
	GetEntities() []domain.AssetData
	UpdateEntity(ar domain.AssetData) error
	GetEntityById(ip string) (int, error)
}

type LogMq interface {
	Produce(log grpcPb.LogItem) error
}

type Service struct {
	Autherization
	CRUDList
	Session
	LogMq
}

func NewService(repos *repository.Repository, hash *hasher.Hash, logger *logger.Logger, logMq LogMq) *Service {
	return &Service{
		Autherization: NewAuthService(repos.Authorization, hash),
		CRUDList:      NewServiceCRUD(repos.CRUDList),
		Session:       newSessionRepo(repos.Session, repos.Authorization, hash, logger),
		LogMq:         logMq,
	}
}
