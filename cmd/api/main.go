package main

import (
	"fmt"
	"os"
	"share-notes-app/internal/container"
	"share-notes-app/internal/routes"
	"share-notes-app/pkg/database"
	"share-notes-app/pkg/utils/config"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func main() {
	// set logger
	logger := logrus.New()

	// load config
	config, err := config.LoadConfig(logger)
	if err != nil {
		logger.Fatal(err)
	} 

	// koneksi ke database
	db, err := database.GetDBConnection(config, logger)
	if err != nil {
		logger.Fatal(err)
	}

	logger.Info(os.Getenv("APP_DATABASE_PASSWORD"), "ini config smtp")
	logger.Info("Berhasil terhubung ke database")

	// buat migrasi
	err = database.CreateMigrationTable(db)
	if err != nil {
		logger.Fatal(err)
	}
	logger.Info("Berhasil membuat migrasi tabel")


	// insialisasi rute
	r := gin.Default()

	container := container.NewContainer(db, config)
	routes.RegisterRoutes(r, container)

	logger.Info(config.SMTP, "ini config smtp")
	logger.Info(os.Getenv("APP_SMTP_AUTH_EMAIL"), "ini config smtp")
	logger.Info(os.Getenv("APP_SMTP_PASSWORD_EMAIL"), "ini config smtp")

	// jalankan server
	PORT := fmt.Sprintf(":%v", config.Server.Port)

	r.Run(PORT)
}