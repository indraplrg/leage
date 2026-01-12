package repositories

import (
	"context"
	"errors"
	"share-notes-app/internal/models"

	"gorm.io/gorm"
)

type AuthenticationRepository interface {
	CreateOne(ctx context.Context, filter any) error
	FindOne(ctx context.Context, filter map[string]any) (*models.User, error)
	DeleteOne(ctx context.Context, id string) error
	GetToken(ctx context.Context, token string) (*models.EmailVerification ,error)
	UpdateOneUsers(ctx context.Context, emailVerify *models.EmailVerification) error
	FindRefreshToken(ctx context.Context, filter map[string]any) (*models.Token, error)
	UpdateRefreshToken(ctx context.Context, token *models.Token) error
}

type authenticationRepository struct {
	db *gorm.DB
}

func NewAuthenticationRepository(db *gorm.DB) AuthenticationRepository {
	return &authenticationRepository{db:db}
}

// Global
func (r *authenticationRepository) CreateOne(ctx context.Context, filter any) error {
	return r.db.WithContext(ctx).Create(filter).Error
}

// User
func (r *authenticationRepository) FindOne(ctx context.Context, filter map[string]any) (*models.User ,error) {
	var User models.User

	tx := r.db.WithContext(ctx)

	for field, value := range filter {
		tx = tx.Where(field+" = ?", value)
	}

	if err := tx.First(&User).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	return &User, nil
}

// Email verification
func (r *authenticationRepository) GetToken(ctx context.Context, token string) (*models.EmailVerification ,error)  {
	var EmailVerification models.EmailVerification 

	if err := r.db.WithContext(ctx).Where("token = ?", token).First(&EmailVerification).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	return &EmailVerification, nil
}

func (r *authenticationRepository) UpdateOneUsers(ctx context.Context, emailVerify *models.EmailVerification) error {
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

// RefreshToken
func (r *authenticationRepository) FindRefreshToken(ctx context.Context, filter map[string]any) (*models.Token, error) {
	var token models.Token
	tx := r.db.WithContext(ctx)

	for field, value := range filter {
		tx = tx.Where(field+" = ?", value)
	}

	if err := tx.First(&token).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}

		return nil, err
	}

	return &token, nil
}

func (r *authenticationRepository) UpdateRefreshToken(ctx context.Context, token *models.Token) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("user_id", token.UserID).Delete(&models.Token{}).Error; err != nil {
			return err
		} 
		
		if err := tx.Create(token).Error; err != nil {
			return err
		}
		return nil
	})
}

func (r *authenticationRepository) DeleteOne(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Where("user_id = ?", id).Delete(&models.Token{}).Error
}

