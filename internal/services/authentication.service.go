package services

import (
	"context"
	"errors"
	"share-notes-app/internal/dtos"
	"share-notes-app/internal/models"
	"share-notes-app/internal/repositories"
	"share-notes-app/pkg/auth"
	"share-notes-app/pkg/mailer"
	"share-notes-app/pkg/token"
	"time"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

type AuthenticationService interface {
	Register(ctx context.Context, dto dtos.UserRequest) (*models.User, error)
	Login(ctx context.Context, dto dtos.LoginRequest) (string, error)
	Logout(ctx context.Context, data *dtos.AuthPayload) error
}

type authenticationService struct {
	repo repositories.AuthenticationRepository
	mailer *mailer.Mailer
}

func NewAuthencticationService(repo repositories.AuthenticationRepository, mailer *mailer.Mailer) AuthenticationService {
	return &authenticationService{repo:repo, mailer:mailer}
}

func (s *authenticationService) Register(ctx context.Context, dto dtos.UserRequest) (*models.User, error) {
	
	// cek email kalau terdaftar
	existing, err := s.repo.FindOne(ctx, map[string]any{
		"email" : dto.Email,
	})

	if err != nil {
		logrus.Info(err)
		return nil, errors.New("gagal mengecek user")
	}

	if existing != nil {
		return nil, errors.New("Email sudah terdaftar!")
	}

	// hashing password
	hashedPassword, err := auth.HashingPassword(dto.Password)
	if err != nil {
		logrus.Info(err)
		return nil, errors.New("gagal hashing password")
	}

	// buat akun 
	User := &models.User{
		Email: dto.Email,
		Username: dto.Username,
		Password: string(hashedPassword),
	}

	err = s.repo.CreateOne(ctx, User)
	if err != nil {
		logrus.Info(err)
		return nil, errors.New("gagal membuat akun")
	}


	// buat verifikasi email
	emailVerification := &models.EmailVerification{
		UserId: User.ID,
		Token: uuid.NewString(),
		ExpiresAt: time.Now().Add(24 * time.Hour),
		CreatedAt: time.Now(),
	}

	err = s.repo.CreateOne(ctx, emailVerification)
	if err != nil {
		logrus.Info(err)
		return nil, errors.New("Gagal membuat verifikasi token")
	}


	// kirim verifikasi email
	if err := s.mailer.SendVerification(dto.Email, emailVerification.Token); err != nil {
		logrus.Info(err)
		return nil, errors.New("gagal mengirim email verifikasi")
	}

	return User, nil
}

func (s *authenticationService) Login(ctx context.Context, dto dtos.LoginRequest) (string, error) {

	// cek akun kalau terdaftar
	user, err := s.repo.FindOne(ctx, map[string]any{
		"username" : dto.Username,
	})

	if err != nil {
		logrus.Info(err)
		return "" ,	errors.New("gagal mengecek user")
	}

	if user == nil {
		return "" , errors.New("account not found!")
	}

	if !user.IsVerified {
		return "" , errors.New("akun belum ter-verifikasi")
	}

	// cek password kalau sama
	err = auth.ComparePassword(user.Password, dto.Password)
	if err != nil {
		logrus.Info(err)
		return "" , errors.New("password yang anda masukkan salah")
	}

	// buat token paseto
	acessToken, err := token.CreateToken(user, time.Now().Add(30 * time.Minute))
	if err != nil {
		logrus.Info(err)
		return "" , errors.New("gagal membuat token")
	}

	refreshToken, err := token.CreateToken(user, time.Now().Add(168 * time.Minute))
	if err != nil {
		logrus.Info(err)
		return "" , errors.New("gagal membuat token")
	}

	// menyimpan token
	Token := models.Token{
		Token: refreshToken,
		UserID: user.ID,
		ExpiredAt: time.Now().Add(168 * time.Minute),
		CreatedAt: time.Now(),
	}

	err = s.repo.CreateOne(ctx, Token)
	if err != nil {
		return "", errors.New("gagal menyimpan token")
	}


	return acessToken, nil
}

func (s *authenticationService) Logout(ctx context.Context, data *dtos.AuthPayload) error {

	userID := data.UserID

	err := s.repo.DeleteOne(ctx, userID)
	if err != nil {
		logrus.Info(err)
		return errors.New("gagal menghapus token")
	}	
	return nil
}