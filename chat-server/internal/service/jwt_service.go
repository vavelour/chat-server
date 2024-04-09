package service

import (
	"errors"
	"github.com/dgrijalva/jwt-go"
	"github.com/vavelour/chat/internal/domain/entities"
	"github.com/vavelour/chat/internal/service/tokens"
)

var (
	errIncorrectType = errors.New("incorrect type")
	errInvalidToken  = errors.New("invalid jwt token")
)

type AuthJWTRepository interface {
	InsertUser(username, password string) error
	GetUser(username string) (entities.User, error)
}

type JwtService struct {
	repos AuthJWTRepository
}

func NewJWTService(r AuthJWTRepository) *JwtService {
	return &JwtService{repos: r}
}

func (s *JwtService) CreateUser(username, password string) (string, error) {
	if err := s.repos.InsertUser(username, password); err != nil {
		return "", err
	}

	return tokens.GenerateToken(username)
}

func (s *JwtService) UserIdentity(token interface{}) (string, error) {
	tkn, ok := token.(string)
	if !ok {
		return "", errIncorrectType
	}

	parsedToken, err := jwt.ParseWithClaims(tkn, &jwt.StandardClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(tokens.SigningKey), nil
	})
	if err != nil {
		return "", err
	}

	if claims, ok := parsedToken.Claims.(*jwt.StandardClaims); ok && parsedToken.Valid {
		return claims.Subject, nil
	}

	return "", errInvalidToken

}
