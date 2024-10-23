package crypto

import (
	"crypto/ed25519"
	"encoding/hex"
	"testing"
)

func TestGenerateEd25519Keys(t *testing.T) {
	// Generate Ed25519 key pair
	publicKeyHex, privateKeyHex, err := generateEd25519Keys()
	if err != nil {
		t.Errorf("Error generating key pair: %v", err)
		return
	}

	// Decode hexadecimal strings to byte slices
	publicKeyBytes, err := hex.DecodeString(publicKeyHex)
	if err != nil {
		t.Errorf("Error decoding public key: %v", err)
		return
	}

	privateKeyBytes, err := hex.DecodeString(privateKeyHex)
	if err != nil {
		t.Errorf("Error decoding private key: %v", err)
		return
	}
	// Sign and verify a message using the generated keys
	message := []byte("hello")
	signature := ed25519.Sign(privateKeyBytes, message)
	// Verify that keys are valid Ed25519 keys
	if !ed25519.Verify(publicKeyBytes, message, signature) {
		t.Errorf("Failed to verify signature")
		return
	}

	t.Logf("Public Key (hex): %s", publicKeyHex)
	t.Logf("Private Key (hex): %s", privateKeyHex)
}
