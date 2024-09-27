package keys

import (
	"crypto/ed25519"
	"encoding/base64"
)

func GenerateED25519KeyPair() (ed25519.PublicKey, ed25519.PrivateKey, error) {
	return ed25519.GenerateKey(nil)
}

func ConvertED25519PrivateKeyToString(privateKey ed25519.PrivateKey) string {
	return base64.StdEncoding.EncodeToString(privateKey)
}

func ConvertED25519PublicKeyToString(publicKey ed25519.PublicKey) string {
	return base64.StdEncoding.EncodeToString(publicKey)
}

func ConvertED25519PrivateKeyFromString(privateKeyString string) (ed25519.PrivateKey, error) {
	privateKeyBytes, err := base64.StdEncoding.DecodeString(privateKeyString)
	if err != nil {
		return nil, err
	}
	return privateKeyBytes, nil
}

func ConvertED25519PublicKeyFromString(publicKeyString string) (ed25519.PublicKey, error) {
	publicKeyBytes, err := base64.StdEncoding.DecodeString(publicKeyString)
	if err != nil {
		return nil, err
	}
	return publicKeyBytes, nil
}
