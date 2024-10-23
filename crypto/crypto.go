package crypto

import (
	"crypto/ed25519"
	"encoding/hex"
)

// generateEd25519Keys generates Ed25519 public-private key pair.
func generateEd25519Keys() (publicKeyHex, privateKeyHex string, err error) {
	// Generate key pair
	pubKey, privKey, err := ed25519.GenerateKey(nil)
	if err != nil {
		return "", "", err
	}

	// Convert keys to hexadecimal format
	publicKeyHex = hex.EncodeToString(pubKey)
	privateKeyHex = hex.EncodeToString(privKey)

	return publicKeyHex, privateKeyHex, nil
}
