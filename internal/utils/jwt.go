package utils

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
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

func DecodePrivateKey(keyData string) (*rsa.PrivateKey, error) {
	block, _ := pem.Decode([]byte(keyData))
	if block == nil {
		return nil, fmt.Errorf("failed to decode private key")
	}

	return x509.ParsePKCS1PrivateKey(block.Bytes)
}
