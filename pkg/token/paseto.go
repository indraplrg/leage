package token

import (
	"os"
	"share-notes-app/internal/models"
	"time"

	"aidanwoods.dev/go-paseto"
	"github.com/google/uuid"
)

func CreateToken(data *models.User, exp time.Time) (string, error) {
	token := paseto.NewToken()
	
	token.SetIssuer("leage")
	token.SetSubject(data.ID.String())
	token.SetString("username", data.Username)
	token.SetString("user_id", data.ID.String())
	token.SetJti(uuid.NewString())
	token.SetIssuedAt(time.Now())
	token.SetNotBefore(time.Now())
	token.SetExpiration(exp)


	secretKey, err := paseto.NewV4AsymmetricSecretKeyFromHex(os.Getenv("APP_PASETO_SECRET_KEY"))
	if err != nil {
		return "", err
	}
	
	signed := token.V4Sign(secretKey, nil)

	return signed, nil
}