package middleware

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func AuthMiddleware(publicKey string) gin.HandlerFunc {
	return func(c *gin.Context) {
		fmt.Println("this is public key: %v", publicKey)
		// Obtém o token do header Authorization
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "authorization header is required"})
			return
		}

		// Remove o prefixo "Bearer "
		tokenString := strings.Replace(authHeader, "Bearer ", "", 1)

		// Parse e valida o token
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			// Verifica se o método de assinatura é o esperado
			if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}

			// Decodifica a chave pública PEM
			key, err := jwt.ParseRSAPublicKeyFromPEM([]byte(publicKey))
			if err != nil {
				return nil, fmt.Errorf("error parsing public key: %v", err)
			}

			return key, nil
		})

		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
			return
		}

		if !token.Valid {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
			return
		}

		c.Next()
	}
}
