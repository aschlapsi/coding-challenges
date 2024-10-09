package crypto

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
)

// Signer defines a contract for different types of signing implementations.
type Signer interface {
	Sign(dataToBeSigned []byte) ([]byte, error)
}

// TODO: implement RSA and ECDSA signing ...

type RSASigner struct {
	KeyPair RSAKeyPair
}

func NewRSASigner(keyPair RSAKeyPair) *RSASigner {
	return &RSASigner{
		KeyPair: keyPair,
	}
}

func (s *RSASigner) Sign(dataToBeSigned []byte) ([]byte, error) {
	hash := sha256.New()
	_, err := hash.Write(dataToBeSigned)
	if err != nil {
		return nil, err
	}
	hashed := hash.Sum(nil)
	signature, err := rsa.SignPSS(rand.Reader, s.KeyPair.Private, crypto.SHA256, hashed, nil)
	if err != nil {
		return nil, err
	}
	return signature, nil
}
