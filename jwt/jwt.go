package jwt

import (
	"fmt"
	"go-mongodb-auth/utils"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type Jwt struct {
	secret string
}

var jt Jwt

func init() {
	err := utils.LoadEnv()
	if err != nil {

	}

	jt = Jwt{secret: os.Getenv("SECRET_KEY")}
}

func CreateToken(email string) (string, *Claims, error) {
	claims, err := NewClaims(email)
	if err != nil {
		return "", &Claims{}, err
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(jt.secret))
	if err != nil {
		return "", &Claims{}, err
	}

	return tokenString, claims, nil
}

func VerifyToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("invalid token signing method")
		}
		return []byte(jt.secret), nil
	})

	if err != nil {
		return nil, fmt.Errorf("error parsing token: %w", err)
	}

	claims, ok := token.Claims.(*Claims)
	if !ok {
		return nil, fmt.Errorf("invalid token claims")
	}

	if time.Now().After(claims.ExpiresAt.Time) {
		return nil, fmt.Errorf("token expired")
	}

	return claims, nil
}
