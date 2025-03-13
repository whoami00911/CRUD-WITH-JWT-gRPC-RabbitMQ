package repository

import (
	"database/sql"
	"webPractice1/internal/domain"
	"webPractice1/pkg/logger"
)

type Session interface {
	CreateRToken(token domain.RefreshSession)
	GetRToken(token string) (domain.RefreshSession, error)
}

type Authorization interface {
	CreateUser(user domain.User) (int, error)
	GetUser(user, password string) (int, error)
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

type Repository struct {
	Authorization
	CRUDList
	Session
}

func NewRepository(db *sql.DB, log *logger.Logger) *Repository {
	return &Repository{
		Authorization: NewAuthUserDbInicialize(db, log),
		CRUDList:      NewCrudDbInicialize(db, log),
		Session:       NewSessionDb(db, log),
	}
}
