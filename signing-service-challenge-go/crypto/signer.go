package crypto

import (
	"crypto"
	"crypto/ecdsa"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"errors"
)

const ALGORITHM_RSA = "RSA"
const ALGORITHM_ECC = "ECC"

var ErrUnknownAlgorithm = errors.New("unknown algorithm")

// Signer defines a contract for different types of signing implementations.
type Signer interface {
	Sign(dataToBeSigned []byte) ([]byte, error)
}

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

type ECDSASigner struct {
	KeyPair ECCKeyPair
}

func NewECDSASigner(keyPair ECCKeyPair) *ECDSASigner {
	return &ECDSASigner{
		KeyPair: keyPair,
	}
}

func (s *ECDSASigner) Sign(dataToBeSigned []byte) ([]byte, error) {
	hash := sha256.New()
	_, err := hash.Write(dataToBeSigned)
	if err != nil {
		return nil, err
	}
	hashed := hash.Sum(nil)
	signature, err := ecdsa.SignASN1(rand.Reader, s.KeyPair.Private, hashed)
	if err != nil {
		return nil, err
	}
	return signature, nil
}

func CreateSigner(algorithm string) (Signer, error) {
	switch algorithm {
	case ALGORITHM_RSA:
		keyGenerator := RSAGenerator{}
		keyPair, err := keyGenerator.Generate()
		if err != nil {
			return nil, err
		}
		return NewRSASigner(*keyPair), nil
	case ALGORITHM_ECC:
		keyGenerator := ECCGenerator{}
		keyPair, err := keyGenerator.Generate()
		if err != nil {
			return nil, err
		}
		return NewECDSASigner(*keyPair), nil
	default:
		return nil, ErrUnknownAlgorithm
	}
}
