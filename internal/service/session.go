package service

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"time"
	"webPractice1/internal/domain"
	"webPractice1/internal/repository"
	"webPractice1/pkg/hasher"
	"webPractice1/pkg/logger"

	"github.com/golang-jwt/jwt"
	"github.com/spf13/viper"
)

type updateSessionRepo struct {
	updateSession repository.Session
	auth          repository.Authorization
	hash          *hasher.Hash
	logger        *logger.Logger
}

func newSessionRepo(repoSession repository.Session, repoAuth repository.Authorization, hash *hasher.Hash, logger *logger.Logger) *updateSessionRepo {
	return &updateSessionRepo{
		updateSession: repoSession,
		auth:          repoAuth,
		hash:          hash,
		logger:        logger,
	}
}

func (usr *updateSessionRepo) CreateRToken(token domain.RefreshSession) {
	usr.updateSession.CreateRToken(token)
}
func (usr *updateSessionRepo) GetRToken(token string) (domain.RefreshSession, error) {
	return usr.updateSession.GetRToken(token)
}

func (usr *updateSessionRepo) GenTokens(user, password string) (string, string, error) {
	//хешировать здесь
	password = usr.hash.GenHashPass(password)
	id := usr.auth.GetUser(user, password)
	if id == 0 {
		return "", "", domain.ErrUserNotFound
	}
	ttl, err := time.ParseDuration(viper.GetString("token.token_ttl"))
	refreshTtl, err := time.ParseDuration(viper.GetString("token.refreshToken_ttl"))
	if err != nil {
		usr.logger.Error(fmt.Sprintf("error time parse: %s", err))
		return "", "", err
	}
	//создать токен здесь
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &jwtToken{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(ttl).Unix(),
			IssuedAt:  time.Now().Unix(),
		},
		UserId: id,
	})
	tokenString, err := token.SignedString([]byte(viper.GetString("token.token_key")))
	if err != nil {
		usr.logger.Error(fmt.Sprintf("failed to get token signed string: %s", err))
		return "", "", err
	}
	refreshToken, err := usr.refreshTokenGen()
	if err != nil {
		usr.logger.Error(fmt.Sprintf("failed to get refresh token: %s", err))
		return "", "", err
	}
	//добавить в бд рефрештокен здесь
	usr.updateSession.CreateRToken(domain.RefreshSession{
		UserID:    id,
		Token:     refreshToken,
		ExpiresAt: time.Now().Add(refreshTtl).UTC(),
	})
	return tokenString, refreshToken, nil
}

func (usr *updateSessionRepo) ParseToken(accessToken string) (int, error) {
	token, err := jwt.ParseWithClaims(accessToken, &jwtToken{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, domain.ErrSignInMethod
		}

		return []byte(viper.GetString("token.token_key")), nil
	})
	if err != nil {
		return 0, err
	}

	claims, ok := token.Claims.(*jwtToken)
	if !ok {
		return 0, domain.ErrTokenClaims
	}

	return claims.UserId, nil
}

func (usr *updateSessionRepo) refreshTokenGen() (string, error) {
	tokenBytes := make([]byte, 32)
	_, err := rand.Read(tokenBytes)
	if err != nil {
		usr.logger.Error(fmt.Sprintf("Failed to generate refresh token: %s", err))
		return "", domain.ErrTokenGen
	}
	//return fmt.Sprintf("%x", tokenBytes), nil
	return base64.URLEncoding.EncodeToString(tokenBytes), nil
}

func (usr *updateSessionRepo) UpdateTokens(token string) (string, string, error) {
	session, err := usr.updateSession.GetRToken(token)
	if err != nil {
		usr.logger.Error(fmt.Sprintf("failed to get refresh token: %s", err))
		return "", "", err
	}
	if session.ExpiresAt.Before(time.Now().UTC()) {
		usr.logger.Error(fmt.Sprintf("token is absolete: %s", domain.ErrObsoleteToken))
		return "", "", domain.ErrObsoleteToken
	}
	//usr.logger.Info(fmt.Sprintf("Checking token expiry: ExpiresAt=%v, Now=%v", session.ExpiresAt, time.Now().UTC()))
	ttl, err := time.ParseDuration(viper.GetString("token.token_ttl"))
	refreshTtl, err := time.ParseDuration(viper.GetString("token.refreshToken_ttl"))
	if err != nil {
		usr.logger.Error(fmt.Sprintf("error time parse: %s", err))
		return "", "", err
	}

	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, &jwtToken{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(ttl).Unix(),
			IssuedAt:  time.Now().Unix(),
		},
		UserId: session.UserID,
	})

	tokenString, err := jwtToken.SignedString([]byte(viper.GetString("token.token_key")))
	if err != nil {
		usr.logger.Error(fmt.Sprintf("failed to get token signed string: %s", err))
		return "", "", err
	}

	refreshToken, err := usr.refreshTokenGen()
	if err != nil {
		usr.logger.Error(fmt.Sprintf("failed to get refresh token: %s", err))
		return "", "", err
	}

	usr.updateSession.CreateRToken(domain.RefreshSession{
		UserID:    session.UserID,
		Token:     refreshToken,
		ExpiresAt: time.Now().Add(refreshTtl).UTC(),
	})
	return tokenString, refreshToken, nil
}
