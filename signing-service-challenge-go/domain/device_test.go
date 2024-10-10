package domain

import (
	gocrypto "crypto"
	"crypto/rsa"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"testing"

	"github.com/fiskaly/coding-challenges/signing-service-challenge/crypto"
)

func TestSign(t *testing.T) {
	var last_signature *Signature

	rsaGenerator := crypto.RSAGenerator{}
	rsaKeyPair, err := rsaGenerator.Generate()
	if err != nil {
		t.Error("Error while generating RSA key pair, got:", err)
	}
	rsaSigner := crypto.NewRSASigner(*rsaKeyPair)
	device := NewSignatureDevice("id", "label", rsaSigner)

	t.Run("Initial message", func(t *testing.T) {
		data := "test_data_number_1"

		signature, err := device.Sign(data)
		if err != nil {
			t.Error("Error while signing, got:", err)
		}

		expected_secured_data := fmt.Sprintf("0_%s_%s", data, base64.StdEncoding.EncodeToString([]byte(device.Id)))
		if signature.Signed_Data != expected_secured_data {
			t.Error("Expected secured data to be", expected_secured_data, "but got", signature.Signed_Data)
		}
		if err := verifySignature(signature, rsaKeyPair.Public); err != nil {
			t.Error("Error while verifying, got:", err)
		}
		last_signature = signature
	})

	t.Run("Second message", func(t *testing.T) {
		data := "test_data_number_2"

		signature, err := device.Sign(data)
		if err != nil {
			t.Error("Error while signing, got:", err)
		}

		expected_secured_data := fmt.Sprintf("1_%s_%s", data, last_signature.Signature)
		if signature.Signed_Data != expected_secured_data {
			t.Error("Expected secured data to be", expected_secured_data, "but got", signature.Signed_Data)
		}
		if err := verifySignature(signature, rsaKeyPair.Public); err != nil {
			t.Error("Error while verifying, got:", err)
		}
	})
}

func verifySignature(signature *Signature, public_key *rsa.PublicKey) error {
	signature_to_verify, _ := base64.StdEncoding.DecodeString(signature.Signature)
	msgHashSum := sha256.Sum256([]byte(signature.Signed_Data))
	return rsa.VerifyPSS(public_key, gocrypto.SHA256, msgHashSum[:], signature_to_verify, nil)
}
