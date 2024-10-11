package crypto

import (
	"crypto"
	"crypto/ecdsa"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"errors"
)

// ALGORITHM_RSA is a constant for the RSA algorithm.
const ALGORITHM_RSA = "RSA"

// ALGORITHM_ECC is a constant for the ECC algorithm.
const ALGORITHM_ECC = "ECC"

// ErrUnknownAlgorithm is an error for unknown algorithms.
var ErrUnknownAlgorithm = errors.New("unknown algorithm")

// Signer defines a contract for different types of signing implementations.
type Signer interface {
	Sign(dataToBeSigned []byte) ([]byte, error)
}

// RSASigner is a concrete implementation of the Signer interface for RSA keys.
type RSASigner struct {
	KeyPair RSAKeyPair
}

// NewRSASigner is a factory to instantiate a new RSASigner.
func NewRSASigner(keyPair RSAKeyPair) *RSASigner {
	return &RSASigner{
		KeyPair: keyPair,
	}
}

// Sign signs the given data with the RSA private key.
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

// ECDSASigner is a concrete implementation of the Signer interface for ECC keys.
type ECDSASigner struct {
	KeyPair ECCKeyPair
}

// NewECDSASigner is a factory to instantiate a new ECDSASigner.
func NewECDSASigner(keyPair ECCKeyPair) *ECDSASigner {
	return &ECDSASigner{
		KeyPair: keyPair,
	}
}

// Sign signs the given data with the ECC private key.
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

// CreateSigner is a factory to instantiate a new Signer based on the given algorithm.
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
