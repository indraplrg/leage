package routes

import (
	"share-notes-app/internal/controllers"

	"github.com/gin-gonic/gin"
)

func AuthorizationRoute(r *gin.Engine, controller *controllers.AuthorizationController) {

	group := r.Group("/api/auth")
	{
		group.GET("/verify-email", controller.VerifyEmail)
	}
}