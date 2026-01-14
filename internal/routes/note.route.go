package routes

import (
	"share-notes-app/internal/container"
	"share-notes-app/internal/middleware"

	"github.com/gin-gonic/gin"
)

func NoteRoute(r *gin.Engine, c *container.Container) {
	general := r.Group("/api/note")
	{
		general.GET("/get-all-notes", c.NoteController.GetAllNotes)
		general.GET("/get-one-note/:id", c.NoteController.GetOneNote)
	}
	
	protected := r.Group("/api/note")
	protected.Use(middleware.VerifyToken(c))
	{
		protected.POST("/create-note", c.NoteController.CreateNote)
		protected.GET("/get-user-notes", c.NoteController.GetUserNotes)
	}
}