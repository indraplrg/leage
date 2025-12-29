package container

import (
	"share-notes-app/internal/controllers"
	"share-notes-app/internal/repositories"
	"share-notes-app/internal/services"
	"share-notes-app/pkg/mailer"
	"share-notes-app/pkg/utils/config"

	"gorm.io/gorm"
)

type Container struct {
	AuthController *controllers.AuthenticationController
	AuthorizationController *controllers.AuthorizationController

	DB *gorm.DB
	Config *config.Config
}

func NewContainer(db *gorm.DB, config *config.Config) *Container {
	mailer := mailer.NewMailer(config)

	authRepo := repositories.NewAuthenticationRepository(db)
	authService := services.NewAuthencticationService(authRepo, mailer)
	authController := controllers.NewAuthenticationController(authService)

	authorizationRepo := repositories.NewAuthorizationRepository(db)
	authorizationsService := services.NewAuthorizationsService(authorizationRepo)
	authorizationController := controllers.NewAuthorizationsController(authorizationsService)

	return &Container{
		AuthController: authController,
		AuthorizationController: authorizationController,

		DB: db,
		Config: config,
	}
}