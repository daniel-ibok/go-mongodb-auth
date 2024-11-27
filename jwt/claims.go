package jwt

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type Claims struct {
	Email string `json:"email"`
	jwt.RegisteredClaims
}

func NewClaims(email string) (*Claims, error) {
	tokenID, err := uuid.NewRandom()
	if err != nil {
		return &Claims{}, err
	}

	return &Claims{
		Email: email,
		RegisteredClaims: jwt.RegisteredClaims{
			ID:        tokenID.String(),
			Issuer:    email,
			Subject:   email,
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24 * 12)),
		},
	}, nil
}
