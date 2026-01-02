package middleware

import (
	"net/http"
	"os"
	"share-notes-app/internal/dtos"
	"strings"
	"time"

	"aidanwoods.dev/go-paseto"
	"github.com/gin-gonic/gin"
)

func VerifyToken() gin.HandlerFunc {
	return func (c *gin.Context)  {
		// Ambil header authoriaztion
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, dtos.BaseResponse{
				Success: false,
				Message: "Authorization header required",
			})
			return
		}

		// Ambil part bearer di header
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, dtos.BaseResponse{
				Success: false,
				Message: "invalid authorization format",
			})
			return
		}

		// Ambil token string
		tokenString := parts[1]

		// Ambil public key
		publicKey, err := paseto.NewV4AsymmetricPublicKeyFromHex(os.Getenv("APP_PASETO_PUBLIC_KEY"))
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, dtos.BaseResponse{
				Success: false,
				Message: "internal server error",
			})
			return
		}

		// Validasi Token
		parser := paseto.NewParser()

		parser.AddRule(paseto.IssuedBy("leage"))
		parser.AddRule(paseto.NotExpired())
		parser.AddRule(paseto.ValidAt(time.Now()))

		parsedToken, err := parser.ParseV4Public(publicKey, tokenString, nil)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, dtos.BaseResponse{
				Success: false,
				Message: "invalid token authorization",
			})
			return
		}

		// Ambil Payload
		userID, _ := parsedToken.GetString("user_id")
		username, _ := parsedToken.GetString("username")

		// Set ke context gin
		c.Set("auth", &dtos.AuthPayload{
			UserID: userID,
			Username: username,
		})
		// Lanjut
		c.Next()

	}
}