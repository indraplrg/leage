package routes

import (
	"share-notes-app/internal/container"
	"share-notes-app/internal/middleware"

	"github.com/gin-gonic/gin"
)


func AuthenticationRoute(r *gin.Engine, c *container.Container) {

	general := r.Group("/api/auth")
	{
		general.POST("/register", c.AuthController.Register)
		general.POST("/login", c.AuthController.Login)
		general.POST("/resend-token")
		general.GET("/verify-email/:token", c.AuthController.VerifyEmail)
	}
	
	protected := r.Group("/api/auth")
	protected.Use(middleware.VerifyToken(c))
	{
		protected.GET("/profile", c.AuthController.Profile)
		protected.POST("/logout" ,c.AuthController.Logout)

	}
}