package api

import (
	"bytes"
	"encoding/json"
	"net/http/httptest"
	"testing"

	"github.com/fiskaly/coding-challenges/signing-service-challenge/crypto"
	"github.com/fiskaly/coding-challenges/signing-service-challenge/domain"
)

func TestSignatureDevice(t *testing.T) {
	s := NewServer(":8080")

	t.Run("invalid request", func(t *testing.T) {
		w := httptest.NewRecorder()
		request := httptest.NewRequest("POST", "/api/v0/signature-device", nil)
		s.SignatureDevice(w, request)
		if w.Code != 400 {
			t.Errorf("Expected status code 400, got %d", w.Code)
		}
	})

	t.Run(("create signature device"), func(t *testing.T) {
		w := httptest.NewRecorder()
		requestBody, err := json.Marshal(CreateSignatureDeviceRequest{
			Id:        "123",
			Algorithm: "ECC",
		})
		if err != nil {
			t.Errorf("Error while marshalling request body: %v", err)
		}
		request := httptest.NewRequest("POST", "/api/v0/signature-device", bytes.NewBuffer(requestBody))
		s.SignatureDevice(w, request)
		if w.Code != 201 {
			t.Errorf("Expected status code 201, got %d", w.Code)
		}
		signatureDevice, err := s.deviceRepository.FindById("123")
		if err != nil {
			t.Fatalf("Error while finding signature device: %v", err)
		}
		if signatureDevice.Id != "123" {
			t.Errorf("Expected signature device with id 123, got %s", signatureDevice.Id)
		}
	})
}

func TestSignature(t *testing.T) {
	s := NewServer(":8080")
	signer, _ := crypto.CreateSigner("RSA")
	signatureDevice := domain.NewSignatureDevice("123", "test_device", signer)
	s.deviceRepository.Save(signatureDevice)
	w := httptest.NewRecorder()
	request := httptest.NewRequest("POST", "/api/v0/signature-device/123/signature", nil)

	s.SignData(w, request)

	if w.Code != 200 {
		t.Errorf("Expected status code 200, got %d", w.Code)
	}
}
