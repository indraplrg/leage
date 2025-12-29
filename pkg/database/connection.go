package database

import (
	"fmt"
	"os"
	"share-notes-app/pkg/utils/config"

	"github.com/sirupsen/logrus"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var Db *gorm.DB

func GetDBConnection(config *config.Config, logger *logrus.Logger) (*gorm.DB, error) {
	// setup connection
	connStr := fmt.Sprintf("host=%s user=%s password=%s port=%s dbname=%s sslmode=%s", 
	config.Database.Host, 
	config.Database.User,
	os.Getenv("APP_DATABASE_PASSWORD"),
	config.Database.Port,
	config.Database.DatabaseName,
	config.Database.SslMode,
	)

	db, err := gorm.Open(postgres.Open(connStr), &gorm.Config{})

	if err != nil {
		logger.WithError(err).Error("Failed connect to database")
		return nil, fmt.Errorf("Failed to open Database connection: %v", err)
	}

	return db, nil
}