package persistence

import (
	"errors"
	"sync"

	"github.com/fiskaly/coding-challenges/signing-service-challenge/domain"
)

// ErrDeviceExists is returned when a device already exists in the repository
var ErrDeviceExists = errors.New("device already exists")

// ErrDeviceNotFound is returned when a device is not found in the repository
var ErrDeviceNotFound = errors.New("device not found")

// SignatureDeviceRepository defines the contract for a signature device repository
type SignatureDeviceRepository interface {
	Save(device *domain.SignatureDevice) error
	FindById(id string) (*domain.SignatureDevice, error)
}

// InMemorySignatureDeviceRepository is an in-memory implementation of a signature device repository
type InMemorySignatureDeviceRepository struct {
	devices map[string]*domain.SignatureDevice
	rwmu    sync.RWMutex
}

// NewInMemorySignatureDeviceRepository creates a new in-memory signature device repository
func NewInMemorySignatureDeviceRepository() *InMemorySignatureDeviceRepository {
	return &InMemorySignatureDeviceRepository{
		devices: make(map[string]*domain.SignatureDevice),
	}
}

// Save saves a signature device in the repository
func (r *InMemorySignatureDeviceRepository) Save(device *domain.SignatureDevice) error {
	r.rwmu.Lock()
	defer r.rwmu.Unlock()

	_, ok := r.devices[device.Id]
	if ok {
		return ErrDeviceExists
	}

	r.devices[device.Id] = device
	return nil
}

// FindById finds a signature device by its id in the repository
func (r *InMemorySignatureDeviceRepository) FindById(id string) (*domain.SignatureDevice, error) {
	r.rwmu.RLock()
	defer r.rwmu.RUnlock()

	device, ok := r.devices[id]
	if !ok {
		return nil, ErrDeviceNotFound
	}
	return device, nil
}
