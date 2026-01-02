package database

import (
	"fmt"
	"os"
	"share-notes-app/configs"

	"github.com/sirupsen/logrus"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var Db *gorm.DB

func GetDBConnection(config *configs.Config, logger *logrus.Logger) (*gorm.DB, error) {
	// setup connection
	connStr := fmt.Sprintf("host=%s user=%s password=%s port=%s dbname=%s sslmode=%s", 
	config.Database.Host, 
	os.Getenv("APP_DATABASE_USERNAME"),
	os.Getenv("APP_DATABASE_PASSWORD"),
	config.Database.Port,
	os.Getenv("APP_DATABASE_NAME"),
	config.Database.SslMode,
	)

	db, err := gorm.Open(postgres.Open(connStr), &gorm.Config{})

	if err != nil {
		logger.WithError(err).Error("Failed connect to database")
		return nil, fmt.Errorf("Failed to open Database connection: %v", err)
	}

	return db, nil
}