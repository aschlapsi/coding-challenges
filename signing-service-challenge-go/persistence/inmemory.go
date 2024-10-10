package persistence

import (
	"errors"

	"github.com/fiskaly/coding-challenges/signing-service-challenge/domain"
)

var ErrDeviceExists = errors.New("device already exists")
var ErrDeviceNotFound = errors.New("device not found")

type SignatureDeviceRepository interface {
	Save(device *domain.SignatureDevice) error
	FindById(id string) (*domain.SignatureDevice, error)
}

type InMemorySignatureDeviceRepository struct {
	devices map[string]*domain.SignatureDevice
}

func NewInMemorySignatureDeviceRepository() *InMemorySignatureDeviceRepository {
	return &InMemorySignatureDeviceRepository{
		devices: make(map[string]*domain.SignatureDevice),
	}
}

func (r *InMemorySignatureDeviceRepository) Save(device *domain.SignatureDevice) error {
	_, ok := r.devices[device.Id]
	if ok {
		return ErrDeviceExists
	}

	r.devices[device.Id] = device
	return nil
}

func (r *InMemorySignatureDeviceRepository) FindById(id string) (*domain.SignatureDevice, error) {
	device, ok := r.devices[id]
	if !ok {
		return nil, ErrDeviceNotFound
	}
	return device, nil
}
