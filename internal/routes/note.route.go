package routes

import (
	"share-notes-app/internal/container"
	"share-notes-app/internal/middleware"

	"github.com/gin-gonic/gin"
)

func NoteRoute(r *gin.Engine, c *container.Container) {
	general := r.Group("/api")
	{
		general.GET("/get-all-notes", c.NoteController.GetAllNotes)
	}
	
	protected := r.Group("/api")
	protected.Use(middleware.VerifyToken(c))
	{
		protected.POST("/create-note", c.NoteController.CreateNote)
		protected.GET("/get-user-notes", c.NoteController.GetUserNotes)
	}
}