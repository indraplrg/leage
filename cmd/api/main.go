package main

import (
	"fmt"
	"share-notes-app/internal/container"
	"share-notes-app/internal/middleware"
	"share-notes-app/internal/routes"
	"share-notes-app/pkg/cache"
	"share-notes-app/pkg/database"
	"share-notes-app/pkg/viper"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func main() {

	// load config
	config, err := viper.LoadConfig()
	if err != nil {
		logrus.Info("gagal me-load config")
		logrus.Fatal(err)
	} 
	logrus.Info("berhasi membaca config")

	// koneksi ke database
	db, err := database.GetDBConnection(config)
	if err != nil {
		logrus.Info("gagal terhubung ke database")
		logrus.Fatal(err)
	}
	logrus.Info("Berhasil terhubung ke database")

	// koneksi ke redis
	redisClient, err := cache.GetValkeyConnection()
	if err != nil {
		logrus.Info("gagal terhubung ke valkey")
		logrus.Fatal(err)
	}
	logrus.Info("Berhasil terhubung ke valkey")

	// buat migrasi
	err = database.CreateMigrationTable(db)
	if err != nil {
		logrus.Info("gagal membuat migrasi")
		logrus.Fatal(err)
	}
	logrus.Info("Berhasil membuat migrasi tabel")


	// insialisasi rute
	r := gin.Default()

	// gunakan middleware
	r.Use(middleware.RateLimit(redisClient))
	r.Use(middleware.RequestLogger())
	r.Use(middleware.Origin())

	container := container.NewContainer(db, config)
	routes.RegisterRoutes(r, container)


	// jalankan server
	PORT := fmt.Sprintf(":%v", config.Server.Port)

	r.Run(PORT)
}