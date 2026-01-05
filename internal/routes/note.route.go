package routes

import (
	"share-notes-app/internal/controllers"
	"share-notes-app/internal/middleware"

	"github.com/gin-gonic/gin"
)

func NoteRoute(r *gin.Engine, controller *controllers.NoteController) {
	group := r.Group("/api")

	{
		group.POST("/create-note", middleware.VerifyToken() ,controller.CreateNote)
	}
}