package utils

import (
	"crypto/rsa"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt"
)

func GenerateJWT(appID int64, privateKey *rsa.PrivateKey) (string, error) {
	now := time.Now().Unix()
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, jwt.StandardClaims{
		IssuedAt:  now,
		ExpiresAt: now + (10 * 60), // JWT valid for 10 minutes
		Issuer:    fmt.Sprintf("%d", appID),
	})

	return token.SignedString(privateKey)
}
