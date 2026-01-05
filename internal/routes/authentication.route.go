package routes

import (
	"share-notes-app/internal/controllers"
	"share-notes-app/internal/middleware"

	"github.com/gin-gonic/gin"
)


func AuthenticationRoute(r *gin.Engine, controller *controllers.AuthenticationController) {

	group := r.Group("/api/auth")
	{
		group.POST("/register", controller.Register)
		group.POST("/login", controller.Login)
		group.POST("/logout", middleware.VerifyToken() ,controller.Logout)
		group.POST("/resend-token")
	}
}