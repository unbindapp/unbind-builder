package utils

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"os"
)

func LoadPrivateKeyFromFile(filePath string) (*rsa.PrivateKey, error) {
	// Load private key
	keyData, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read private key: %v", err)
	}

	block, _ := pem.Decode(keyData)
	if block == nil {
		return nil, fmt.Errorf("failed to decode private key")
	}

	return x509.ParsePKCS1PrivateKey(block.Bytes)
}
