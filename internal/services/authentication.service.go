package services

import (
	"context"
	"fmt"
	"share-notes-app/internal/dtos"
	"share-notes-app/internal/models"
	"share-notes-app/internal/repositories"
	"share-notes-app/pkg/auth"
	"share-notes-app/pkg/mailer"
	"time"

	"github.com/google/uuid"
)

type AuthenticationService interface {
	Register(ctx context.Context, dto dtos.UserRequest) (string, error)
}

type authenticationService struct {
	repo repositories.AuthenticationRepository
	mailer *mailer.Mailer
}

func NewAuthencticationService(repo repositories.AuthenticationRepository, mailer *mailer.Mailer) AuthenticationService {
	return &authenticationService{repo:repo, mailer:mailer}
}

func (s *authenticationService) Register(ctx context.Context, dto dtos.UserRequest) (string, error) {
	
	// cek email kalau terdaftar
	existing, err := s.repo.FindByEmail(ctx, dto.Email)

	if err != nil {
		return "", fmt.Errorf("gagal mengecek user: %w", err)
	}

	if existing != nil {
		return "", fmt.Errorf("Email sudah terdaftar!")
	}

	// hashing password
	hashedPassword, err := auth.HashingPassword(dto.Password)
	if err != nil {
		return "", fmt.Errorf("gagal hashing password: %w", err)
	}

	// buat akun 
	User := &models.User{
		Email: dto.Email,
		Username: dto.Username,
		Password: string(hashedPassword),
	}

	err = s.repo.CreateUser(ctx, User)
	if err != nil {
		return "", fmt.Errorf("gagal membuat akun: %w", err)
	}


	// buat verifikasi email
	emailVerification := &models.EmailVerification{
		UserId: User.ID,
		Token: uuid.NewString(),
		ExpiresAt: time.Now().Add(24 * time.Hour),
		CreatedAt: time.Now(),
	}

	err = s.repo.CreateEmailVerification(ctx, emailVerification)
	if err != nil {
		return "", fmt.Errorf("Gagal membuat verifikasi token: %w", err)
	}


	// kirim verifikasi email
	if err := s.mailer.SendVerification(dto.Email, emailVerification.Token); err != nil {
		return "", fmt.Errorf("gagal mengirim email verifikasi: %w", err)
	}

	return "berhasil mengirim verifikasi ke email anda", nil
}