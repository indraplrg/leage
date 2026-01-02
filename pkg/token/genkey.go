package token

import (
	"fmt"

	"aidanwoods.dev/go-paseto"
)

func generateKey() {
	secretKey := paseto.NewV4AsymmetricSecretKey()
	publicKey := secretKey.Public()

	fmt.Println(secretKey.ExportHex(), "ini secret key")
	fmt.Println(publicKey.ExportHex(), "ini public key")
}