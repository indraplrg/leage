package services

import (
	"context"
	"errors"
	"share-notes-app/internal/dtos"
	"share-notes-app/internal/models"
	"share-notes-app/internal/repositories"
	"share-notes-app/pkg/apperror"
	"share-notes-app/pkg/auth"
	"share-notes-app/pkg/mailer"
	"share-notes-app/pkg/token"
	"time"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

type AuthenticationService interface {
	Register(ctx context.Context, dto dtos.UserRegisterRequest) (*models.User, error)
	Login(ctx context.Context, dto dtos.UserLoginRequest) (*dtos.LoginData, error)
	Profile(ctx context.Context, payload *dtos.AuthPayload) (*models.User, error)
	Logout(ctx context.Context, data *dtos.AuthPayload) (error)
	VerifyEmail(ctx context.Context, token string) (string, error)
	ValidateRefreshToken(ctx context.Context, userID, refreshToken string) (bool, error)
}

type authenticationService struct {
	repo repositories.AuthenticationRepository
	mailer *mailer.Mailer
}

func NewAuthencticationService(repo repositories.AuthenticationRepository, mailer *mailer.Mailer) AuthenticationService {
	return &authenticationService{repo:repo, mailer:mailer}
}

func (s *authenticationService) Register(ctx context.Context, dto dtos.UserRegisterRequest) (*models.User, error) {
	
	// cek email kalau terdaftar
	existing, err := s.repo.FindOne(ctx, map[string]any{
		"email" : dto.Email,
	})

	if err != nil {
		logrus.WithError(err)
		return nil, errors.New("gagal mengecek user")
	}

	if existing != nil {
		return nil, errors.New("Email sudah terdaftar!")
	}

	// hashing password
	hashedPassword, err := auth.HashingPassword(dto.Password)
	if err != nil {
		logrus.WithError(err)
		return nil, errors.New("gagal hashing password")
	}

	// buat akun 
	User := &models.User{
		Email: dto.Email,
		Username: dto.Username,
		Password: string(hashedPassword),
		IsVerified: false,
	}

	err = s.repo.CreateOne(ctx, User)
	if err != nil {
		logrus.WithError(err)
		return nil, errors.New("gagal membuat akun")
	}


	// buat verifikasi email
	emailVerification := &models.EmailVerification{
		UserID: User.ID,
		Token: uuid.NewString(),
		ExpiresAt: time.Now().Add(24 * time.Hour),
		CreatedAt: time.Now(),
		IsUsed: false,
	}

	err = s.repo.CreateOne(ctx, emailVerification)
	if err != nil {
		logrus.WithError(err)
		return nil, errors.New("Gagal membuat verifikasi token")
	}


	// kirim verifikasi email
	if err := s.mailer.SendVerification(dto.Email, emailVerification.Token); err != nil {
		logrus.WithError(err)
		return nil, errors.New("gagal mengirim email verifikasi")
	}

	return User, nil
}

func (s *authenticationService) Login(ctx context.Context, dto dtos.UserLoginRequest) (*dtos.LoginData, error) {

	// cek akun kalau terdaftar
	user, err := s.repo.FindOne(ctx, map[string]any{
		"username" : dto.Username,
	})

	if err != nil {
		logrus.WithError(err)
		return nil ,	errors.New("gagal mengecek user")
	}

	if user == nil {
		return nil , errors.New("account not found!")
	}

	if !user.IsVerified {
		return nil , errors.New("akun belum ter-verifikasi")
	}

	// cek password kalau sama
	err = auth.ComparePassword(user.Password, dto.Password)
	if err != nil {
		logrus.WithError(err)
		return nil , errors.New("password yang anda masukkan salah")
	}

	// buat token paseto
	acessToken, err := token.CreateToken(user.Username, user.ID.String(), time.Now().Add(30 * time.Minute))
	if err != nil {
		logrus.WithError(err)
		return nil , errors.New("gagal membuat token")
	}

	refreshToken, err := token.CreateToken(user.Username, user.ID.String(), time.Now().Add(168 * time.Hour))
	if err != nil {
		logrus.WithError(err)
		return nil , errors.New("gagal membuat token")
	}

	// Hashing refresh token
	hashedRefreshToken := auth.HasingRefreshToken(refreshToken)

	// membuat model token
	Token := &models.Token{
		Token: string(hashedRefreshToken),
		UserID: user.ID,
		ExpiredAt: time.Now().Add(168 * time.Hour),
		CreatedAt: time.Now(),
	}

	// menyimpan token
	err = s.repo.UpdateRefreshToken(ctx, Token)
	if err != nil {
		logrus.WithError(err)
		return nil, errors.New("gagal menyimpan token")
	}


	return &dtos.LoginData{AccessToken: acessToken, RefreshToken: refreshToken}, nil
}

func (s *authenticationService) Profile(ctx context.Context, payload *dtos.AuthPayload) (*models.User, error) {
	userID := payload.UserID
	user, err := s.repo.FindOne(ctx, map[string]any{
		"id" : userID,
	})

	if err != nil {
		logrus.WithError(err)
		return nil, apperror.Failed("check User")
	}

	if user == nil {
		return nil, apperror.NotFound("User")
	} 

	return user, nil
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

func (s *authenticationService) VerifyEmail(ctx context.Context, token string) (string, error) {

	// cek token kalau ada
	emailVerify, err := s.repo.GetToken(ctx, token)

	if err != nil {
		logrus.WithError(err)
		return "", errors.New("gagal mengambil token")
	}

	if emailVerify == nil {
		return "", errors.New("token tidak ditemukan")
	}

	if emailVerify.IsUsed {
		return "", errors.New("token sudah digunakan")
	}

	if time.Now().After(emailVerify.ExpiresAt) {
		return "", errors.New("token sudah kadaluarsa")
	}


	// update status verifikasi
	err = s.repo.UpdateOneUsers(ctx, emailVerify)
	if err != nil {
		logrus.WithError(err)
		return "", errors.New("gagal mengupdate status")
	} 

	return "berhasil meng-update status verifikasi", nil
}

func (s *authenticationService) ValidateRefreshToken(ctx context.Context, userID, refreshToken string) (bool, error) {
	// hash refresh token
	hashedRefreshToken := auth.HasingRefreshToken(refreshToken)

	// ambil refresh token dari db
	tokenFromDB, err := s.repo.FindRefreshToken(ctx, map[string]any{
		"token": hashedRefreshToken,
	})

	// validasi
	if err != nil {
		logrus.WithError(err)
		return false, errors.New("gagal mengambil refresh token")
	}

	if tokenFromDB == nil {
		return false, errors.New("refresh token tidak ditemukan")
	}

	if time.Now().After(tokenFromDB.ExpiredAt) {
		_ = s.repo.DeleteOne(ctx, userID)
		return false, errors.New("refresh token tidak valid")
	}

	return true, nil
}

// func (s *authenticationService) ResendToken(ctx context.Context, data *dtos.ResendTokenRequest) error {
	
// 	// Cek akun kalau terdaftar
// 	existing, err := s.repo.FindOne(ctx, map[string]any{
// 		"email": data.Email,
// 	})

// 	if err != nil {
// 		logrus.WithError(err)
// 		return errors.New("gagal mengecek user")
// 	}

// 	if existing == nil {
// 		return errors.New("Email belum terdaftar!")
// 	}

// 	if existing.IsVerified {
// 		return errors.New("Email telah terverifikasi!")
// 	}

// 	// Ambil token
// 	token, err := s.repo.FindOne(ctx, map[string]any{
// 		"user_id": existing.ID,
// 	})


// 	if err := s.mailer.SendVerification(data.Email, token); err != nil {
// 		logrus.WithError(err)
// 		return errors.New("gagal mengirim email verifikasi")
// 	}
	
// 	return nil
// }

