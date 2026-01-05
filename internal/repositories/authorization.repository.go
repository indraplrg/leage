package repositories

import (
	"context"
	"errors"
	"share-notes-app/internal/models"

	"gorm.io/gorm"
)

type AuthorizationRepositories interface {
	GetToken(ctx context.Context, token string) (*models.EmailVerification ,error)
	UpdateOneUsers(ctx context.Context, emailVerify *models.EmailVerification) error
}

type authorizationRepository struct {
	db *gorm.DB
}

func NewAuthorizationRepository(db *gorm.DB) AuthorizationRepositories {
	return &authorizationRepository{db:db}
}

func (r *authorizationRepository) GetToken(ctx context.Context, token string) (*models.EmailVerification ,error)  {
	var EmailVerification models.EmailVerification 

	if err := r.db.WithContext(ctx).Where("token = ?", token).First(&EmailVerification).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	return &EmailVerification, nil
}

func (r *authorizationRepository) UpdateOneUsers(ctx context.Context, emailVerify *models.EmailVerification) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// update user
		if err := tx.Model(&models.User{}).Where("id = ?", emailVerify.UserID).Update("is_verified", true).Error; err != nil {
			return err
		}

		// update token
		if err := tx.Model(&models.EmailVerification{}).Where("id = ?", emailVerify.ID).Update("is_used", true).Error; err != nil {
			return err
		}

		return nil
	})
}