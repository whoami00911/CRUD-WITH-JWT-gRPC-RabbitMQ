package service

import (
	"webPractice1/internal/domain"
	"webPractice1/internal/repository"
	"webPractice1/pkg/hasher"

	"github.com/golang-jwt/jwt"
)

type AuthService struct {
	repo repository.Authorization
	hash *hasher.Hash
}

type jwtToken struct {
	jwt.StandardClaims
	UserId int
}

func NewAuthService(repo repository.Authorization, hash *hasher.Hash) *AuthService {
	return &AuthService{
		repo: repo,
		hash: hash,
	}
}

func (as *AuthService) CreateUser(user domain.User) (int, error) {
	user.Password = as.hash.GenHashPass(user.Password)
	return as.repo.CreateUser(user)
}
