package routes

import (
	"share-notes-app/internal/container"
	"share-notes-app/internal/middleware"

	"github.com/gin-gonic/gin"
)

func NoteRoute(r *gin.Engine, c *container.Container) {
	protected := r.Group("/api")
	protected.Use(middleware.VerifyToken(c))
	{
		protected.POST("/create-note", c.NoteController.CreateNote)
	}
}