package routes

import (
	"share-notes-app/internal/controllers"

	"github.com/gin-gonic/gin"
)


func AuthenticationRoute(r *gin.Engine, controller *controllers.AuthenticationController) {

	group := r.Group("/api/auth")
	{
		group.POST("/register", controller.Register)
	}
}