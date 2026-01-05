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
}

type authenticationRepository struct {
	db *gorm.DB
}

func NewAuthenticationRepository(db *gorm.DB) AuthenticationRepository {
	return &authenticationRepository{db:db}
}

func (r *authenticationRepository) CreateOne(ctx context.Context, filter any) error {
	return r.db.WithContext(ctx).Create(filter).Error
}

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

func (r *authenticationRepository) DeleteOne(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Where("user_id = ?", id).Model(&models.Token{}).Error
}