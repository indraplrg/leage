package repositories

import "gorm.io/gorm"

type AuthorizationRepositories interface {
	
}

type authorizationRepository struct {
	db *gorm.DB
}

func NewAuthorizationRepository(db *gorm.DB) AuthorizationRepositories {
	return &authorizationRepository{db:db}
}