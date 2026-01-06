package token

import (
	"os"
	"time"

	"aidanwoods.dev/go-paseto"
	"github.com/google/uuid"
)

func CreateToken(username, userID string, exp time.Time) (string, error) {
	token := paseto.NewToken()
	
	token.SetIssuer("leage")
	token.SetSubject(userID)
	token.SetString("username", username)
	token.SetString("user_id", userID)
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