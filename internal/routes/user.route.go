package routes

import (
	"share-notes-app/internal/container"
	"share-notes-app/internal/middleware"

	"github.com/gin-gonic/gin"
)


func AuthenticationRoute(r *gin.Engine, c *container.Container) {

	group := r.Group("/api/auth")
	{
		group.POST("/register", c.AuthController.Register)
		group.POST("/login", c.AuthController.Login)
		group.POST("/resend-token")
		group.GET("/verify-email/:token", c.AuthController.VerifyEmail)
	}
	
	protected := r.Group("/api/auth")
	protected.Use(middleware.VerifyToken(c))
	{
		protected.POST("/logout" ,c.AuthController.Logout)

	}
}