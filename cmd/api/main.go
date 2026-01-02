package main

import (
	"fmt"
	"share-notes-app/internal/container"
	"share-notes-app/internal/middleware"
	"share-notes-app/internal/routes"
	"share-notes-app/pkg/cache"
	"share-notes-app/pkg/database"
	"share-notes-app/pkg/utils/viper"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func main() {
	// set logger
	logger := logrus.New()

	// load config
	config, err := viper.LoadConfig(logger)
	if err != nil {
		logger.Fatal(err)
	} 

	// koneksi ke database
	db, err := database.GetDBConnection(config, logger)
	if err != nil {
		logger.Fatal(err)
	}
	logger.Info("Berhasil terhubung ke database")

	// koneksi ke redis
	redisClient, err := cache.GetRedisConnection()
	if err != nil {
		logger.Fatal(err)
	}
	logger.Info("Berhasil terhubung ke valkey")

	// buat migrasi
	err = database.CreateMigrationTable(db)
	if err != nil {
		logger.Fatal(err)
	}
	logger.Info("Berhasil membuat migrasi tabel")


	// insialisasi rute
	r := gin.Default()

	// gunakan middleware
	r.Use(middleware.RateLimit(redisClient))

	container := container.NewContainer(db, config)
	routes.RegisterRoutes(r, container)


	// jalankan server
	PORT := fmt.Sprintf(":%v", config.Server.Port)

	r.Run(PORT)
}