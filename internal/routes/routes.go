package routes

import (
	"share-notes-app/internal/container"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.Engine, c *container.Container) {
	AuthenticationRoute(r, c.AuthController)
	AuthorizationRoute(r, c.AuthorizationController) 
}