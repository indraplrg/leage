package repositories

import (
	"context"
	"errors"
	"share-notes-app/internal/models"

	"gorm.io/gorm"
)

type AuthenticationRepository interface {
	CreateUser(ctx context.Context, user *models.User) error
	CreateEmailVerification(ctx context.Context, verification *models.EmailVerification) error
	FindByEmail(ctx context.Context, email string) (*models.User, error)
}

type authenticationRepository struct {
	db *gorm.DB
}

func NewAuthenticationRepository(db *gorm.DB) AuthenticationRepository {
	return &authenticationRepository{db:db}
}

func (r *authenticationRepository) CreateUser(ctx context.Context, user *models.User) error {
	return r.db.WithContext(ctx).Create(user).Error
}

func (r *authenticationRepository) CreateEmailVerification(ctx context.Context, verification *models.EmailVerification) error {
	return r.db.WithContext(ctx).Create(verification).Error
}

func (r *authenticationRepository) FindByEmail(ctx context.Context, email string) (*models.User ,error) {
	var User models.User

	if err := r.db.WithContext(ctx).Where("email = ?", email).First(&User).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	return &User, nil
}