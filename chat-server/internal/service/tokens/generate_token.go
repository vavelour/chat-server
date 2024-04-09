package tokens

import (
	"github.com/dgrijalva/jwt-go"
	"time"
)

const (
	SigningKey = "valera_super_star_for_real"
	tokenTTL   = 12 * time.Hour
)

func GenerateToken(username string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &jwt.StandardClaims{
		ExpiresAt: time.Now().Add(tokenTTL).Unix(),
		IssuedAt:  time.Now().Unix(),
		Subject:   username,
	})

	return token.SignedString([]byte(SigningKey))
}
