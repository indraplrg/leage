package container

import (
	"share-notes-app/configs"
	"share-notes-app/internal/controllers"
	"share-notes-app/internal/repositories"
	"share-notes-app/internal/services"
	"share-notes-app/pkg/mailer"

	"gorm.io/gorm"
)

type Container struct {
	AuthController *controllers.AuthenticationController
	NoteController *controllers.NoteController

	AuthService services.AuthenticationService

	DB *gorm.DB
	Config *configs.Config
}

func NewContainer(db *gorm.DB, config *configs.Config) *Container {
	mailer := mailer.NewMailer(config)

	authRepo := repositories.NewAuthenticationRepository(db)
	authService := services.NewAuthencticationService(authRepo, mailer)
	authController := controllers.NewAuthenticationController(authService)

	noteRepo := repositories.NewNoteRepository(db)
	noteService := services.NewNoteService(noteRepo)
	noteController := controllers.NewNoteController(noteService)

	return &Container{
		AuthController: authController,
		NoteController: noteController,

		AuthService: authService,

		DB: db,
		Config: config,
	}
}