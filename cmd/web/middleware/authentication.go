package middleware

import (
	"net/http"
	"os"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

var JwtKey = []byte(os.Getenv("JWT_SECRET_KEY"))

type Credentials struct {
	Password string `json:"password"`
	Username string `json:"username"`
}

type Claims struct {
	Username string `json:"username"`
	jwt.StandardClaims
}

func AuthRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		receivedToken := strings.ReplaceAll(c.GetHeader("Authorization"), "Bearer ", "")
		claims := &Claims{}

		token, err := jwt.ParseWithClaims(receivedToken, claims, func(token *jwt.Token) (interface{}, error) {
			return JwtKey, nil
		})
		if err != nil {
			if err == jwt.ErrSignatureInvalid {
				c.AbortWithStatusJSON(http.StatusUnauthorized, err)
				return
			}
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Invalid token"})
			return
		}
		if !token.Valid {
			c.AbortWithStatusJSON(http.StatusUnauthorized, err)
			return
		}
		c.Next()
	}
}
