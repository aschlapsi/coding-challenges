package crypto

import (
	"crypto"
	"crypto/rsa"
	"crypto/sha256"
	"testing"
)

func TestRSASigner_Sign(t *testing.T) {
	rsaGenerator := &RSAGenerator{}
	rsaKeyPair, _ := rsaGenerator.Generate()
	message := []byte("test_data")

	rsaSigner := NewRSASigner(*rsaKeyPair)
	signature, err := rsaSigner.Sign(message)
	if err != nil {
		t.Error("Error while signing, got:", err)
	}

	msgHashSum := sha256.Sum256(message)
	err = rsa.VerifyPSS(rsaKeyPair.Public, crypto.SHA256, msgHashSum[:], signature, nil)
	if err != nil {
		t.Error("Error while verifying, got:", err)
	}
}
