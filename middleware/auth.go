package middleware

import (
	"go-mongodb-auth/jwt"
	"go-mongodb-auth/models"
	"go-mongodb-auth/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {

		// read cookie
		cookie, err := c.Cookie("Authorization")
		if err != nil || cookie == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"status":  http.StatusUnauthorized,
				"message": http.StatusText(http.StatusUnauthorized),
			})
			return
		}

		token := utils.RetrieveToken(cookie)

		// verify token
		claims, err := jwt.VerifyToken(token)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"status":  http.StatusInternalServerError,
				"message": err.Error(),
			})
			return
		}

		// retrieve user from database
		user, err := models.GetUserByEmail(claims.Email)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"status":  http.StatusInternalServerError,
				"message": err.Error(),
			})
			return
		}

		// store new user credentials
		c.Set("user", user)
		c.Next()
	}
}
