package middleware

import (
	"context"
	"net/http"
	"os"
	"share-notes-app/helper"
	"share-notes-app/internal/container"
	"share-notes-app/internal/dtos"
	"share-notes-app/pkg/token"
	"strings"
	"time"

	"aidanwoods.dev/go-paseto"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func VerifyToken(container *container.Container) gin.HandlerFunc {
	return func (c *gin.Context)  {
		// Ambil header authoriaztion
		accesstoken, err := helper.GetCookie(c, "access_paseto_token")
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, dtos.BaseResponse{
				Success: false,
				Message: "access token required",
			})
			return
		}

		// Ambil public key
		publicKey, err := paseto.NewV4AsymmetricPublicKeyFromHex(os.Getenv("APP_PASETO_PUBLIC_KEY"))
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, dtos.BaseResponse{
				Success: false,
				Message: "internal server error",
			})
			return
		}

		parser := paseto.NewParser()
		parser.AddRule(paseto.IssuedBy("leage"))
		
		// Validasi Token
		parsedToken, err := parser.ParseV4Public(publicKey, accesstoken, nil)
		
		if err != nil {
			if strings.Contains(err.Error(), "expired") {
				logrus.Info("token expire bro!!!")
				handleRefresh(c, container)
				return
			}

			c.AbortWithStatusJSON(http.StatusUnauthorized, dtos.BaseResponse{
				Success: false,
				Message: "invalid token",
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

func handleRefresh(c *gin.Context, container *container.Container) {
	refreshToken, err := helper.GetCookie(c,"refresh_paseto_token")
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"message": "refresh token required",
		})
		return
	}

	// Ambil public key
	publicKey, err := paseto.NewV4AsymmetricPublicKeyFromHex(os.Getenv("APP_PASETO_PUBLIC_KEY"))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, dtos.BaseResponse{
			Success: false,
			Message: "internal server error",
		})
		return
	}

	// Validasi refresh token
	parser := paseto.NewParser()
	parser.AddRule(paseto.IssuedBy("leage"))
	parser.AddRule(paseto.NotExpired())
	parser.AddRule(paseto.ValidAt(time.Now()))

	parsedRefreshToken, err := parser.ParseV4Public(publicKey, refreshToken, nil)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, dtos.BaseResponse{
			Success: false,
			Message: "refresh token expired, please login again",
		})
		return
	}

	// ambil payload refresh token
	userID, err := parsedRefreshToken.GetString("user_id")
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, dtos.BaseResponse{
			Success: false,
			Message: "invalid refresh token payload",
		})
		return
	}

	username, err := parsedRefreshToken.GetString("username")
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, dtos.BaseResponse{
			Success: false,
			Message: "invalid refresh token payload",
		})
		return
	}

	// Validasi token kalau sama
	ctx := context.Background()
	ok, err := container.AuthService.ValidateRefreshToken(ctx, userID, refreshToken)
	if err != nil || !ok {
		c.AbortWithStatusJSON(http.StatusUnauthorized, dtos.BaseResponse{
			Success: false,
			Message: "invalid refresh token, please login again",
		})
		return
	}

	// buat access token baru
	newAccessToken, err := token.CreateToken(username, userID, time.Now().Add(30 * time.Minute))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, dtos.BaseResponse{
			Success: false,
			Message: "failed to generate new access token",
		})
		return
	}

	// set access token baru
	helper.SetCookie(c, "access_paseto_token", newAccessToken, 168*3600)

	c.Set("auth", &dtos.AuthPayload{
			UserID: userID,
			Username: username,
		})
	
	// Lanjut
	c.Next()
}