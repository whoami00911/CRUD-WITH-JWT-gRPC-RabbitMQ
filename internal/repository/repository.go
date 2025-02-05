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
	GetUser(user, password string) int
}

type CRUDList interface {
	AddEntity(ar domain.AssetData)
	DeleteAllEntitiesDB()
	DeleteEntityDB(ip string)
	GetEntity(ip string) *domain.AssetData
	GetEntities() []domain.AssetData
	UpdateEntity(ar domain.AssetData)
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
