package utils

import (
	"fmt"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/joho/godotenv"
)

func LoadEnv() error {
	err := godotenv.Load()
	return err
}

func GetUUID() string {
	tokenID, _ := uuid.NewRandom()
	return tokenID.String()
}

func RetrieveToken(cookie string) string {
	data := strings.Split(cookie, " ")
	return data[1]
}

func GenerateAndStoreCookie(c *gin.Context, token string) {
	token = fmt.Sprintf("Bearer %s", token)
	c.SetCookie("Authorization", token, 24*7*3600, "/", "", true, true)
}
