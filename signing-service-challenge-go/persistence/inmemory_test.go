package persistence

import (
	"testing"

	"github.com/fiskaly/coding-challenges/signing-service-challenge/crypto"
	"github.com/fiskaly/coding-challenges/signing-service-challenge/domain"
)

func TestInMemorySignatureDeviceRepository(t *testing.T) {
	repo := NewInMemorySignatureDeviceRepository()
	rsaGenerator := crypto.RSAGenerator{}
	rsaKeyPair, err := rsaGenerator.Generate()
	if err != nil {
		t.Error("Error while generating RSA key pair, got:", err)
	}
	rsaSigner := crypto.NewRSASigner(*rsaKeyPair)
	device := domain.NewSignatureDevice("1", "test_device", rsaSigner)
	repo.Save(device)

	t.Run("Save_DuplicateDevice", func(t *testing.T) {
		otherDevice := domain.NewSignatureDevice("1", "other_device", rsaSigner)
		err := repo.Save(otherDevice)
		if err != ErrDeviceExists {
			t.Error("Expected to get ErrDeviceExists, but got:", err)
		}
	})

	t.Run("FindById_DeviceExists", func(t *testing.T) {
		foundDevice, err := repo.FindById("1")
		if err != nil {
			t.Error("Expected to find device with id 1, but got error:", err)
		}
		if foundDevice == nil {
			t.Error("Expected to find device with id 1, but got nil")
		} else if foundDevice.Id != device.Id {
			t.Error("Expected to find device with id 1, but got device with id", foundDevice.Id)
		}
	})

	t.Run("FindById_DeviceDoesNotExist", func(t *testing.T) {
		foundDevice, err := repo.FindById("2")
		if err == nil {
			t.Error("Expected to not find device with id 2, but got device with id", foundDevice.Id)
		}
	})
}
