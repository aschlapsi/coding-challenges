package api

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/fiskaly/coding-challenges/signing-service-challenge/crypto"
	"github.com/fiskaly/coding-challenges/signing-service-challenge/domain"
	"github.com/google/uuid"
)

// Create a new signature device
func (s *Server) SignatureDevice(response http.ResponseWriter, request *http.Request) {
	switch request.Method {
	case http.MethodGet:
		s.listSignatureDevices(response, request)
	case http.MethodPost:
		s.createSignatureDevice(response, request)
	default:
		WriteErrorResponse(response, http.StatusMethodNotAllowed, []string{
			http.StatusText(http.StatusMethodNotAllowed),
		})
	}
}

type SignatureDevice struct {
	Id    string `json:"id"`
	Label string `json:"label"`
}

func (s *Server) listSignatureDevices(response http.ResponseWriter, _ *http.Request) {
	signatureDevices, err := s.deviceRepository.FindAll()
	if err != nil {
		log.Printf("Error while finding signature devices: %v", err)
		WriteInternalError(response)
		return
	}

	devices := make([]*SignatureDevice, 0)
	for device := range signatureDevices {
		devices = append(devices, &SignatureDevice{
			Id:    signatureDevices[device].Id,
			Label: signatureDevices[device].Label,
		})
	}

	WriteAPIResponse(response, http.StatusOK, devices)
}

type CreateSignatureDeviceRequest struct {
	Id        string `json:"id"`
	Label     string `json:"label"`
	Algorithm string `json:"algorithm"`
}

type CreateSignatureDeviceResponse struct {
	Id        string `json:"id"`
	Label     string `json:"label"`
	Algorithm string `json:"algorithm"`
}

func (s *Server) createSignatureDevice(response http.ResponseWriter, request *http.Request) {
	var createSignatureDeviceRequest CreateSignatureDeviceRequest
	err := json.NewDecoder(request.Body).Decode((&createSignatureDeviceRequest))
	if err != nil {
		log.Printf("Error while decoding request body: %v", err)
		WriteErrorResponse(response, http.StatusBadRequest, []string{
			http.StatusText(http.StatusBadRequest),
		})
		return
	}

	signer, err := crypto.CreateSigner(createSignatureDeviceRequest.Algorithm)
	if err != nil {
		log.Printf("Error while creating signer: %v", err)
		WriteInternalError(response)
		return
	}

	uuid, err := uuid.NewRandom()
	if err != nil {
		log.Printf("Error while generating UUID: %v", err)
		WriteInternalError(response)
		return
	}

	signatureDevice := domain.NewSignatureDevice(
		uuid.String(),
		createSignatureDeviceRequest.Label,
		signer,
	)
	err = s.deviceRepository.Save(signatureDevice)
	if err != nil {
		log.Printf("Error while saving signature device: %v", err)
		WriteErrorResponse(response, http.StatusBadRequest, []string{
			err.Error(),
		})
		return
	}

	createSignatureDeviceResponse := CreateSignatureDeviceResponse{
		Id:        signatureDevice.Id,
		Label:     signatureDevice.Label,
		Algorithm: createSignatureDeviceRequest.Algorithm,
	}
	WriteAPIResponse(response, http.StatusCreated, createSignatureDeviceResponse)
}

type SignDataRequest struct {
	Id   string `json:"id"`
	Data string `json:"data"`
}

type SignDataResponse struct {
	Signature  string `json:"signature"`
	SignedData string `json:"signed_data"`
}

// Sign data with a signature device
func (s *Server) SignData(response http.ResponseWriter, request *http.Request) {
	if request.Method != http.MethodPost {
		WriteErrorResponse(response, http.StatusMethodNotAllowed, []string{
			http.StatusText(http.StatusMethodNotAllowed),
		})
		return
	}

	var signDataRequest SignDataRequest
	err := json.NewDecoder(request.Body).Decode((&signDataRequest))
	if err != nil {
		log.Printf("Error while decoding request body: %v", err)
		WriteErrorResponse(response, http.StatusBadRequest, []string{
			http.StatusText(http.StatusBadRequest),
		})
		return
	}

	signatureDevice, err := s.deviceRepository.FindById(signDataRequest.Id)
	if err != nil {
		log.Printf("Error while finding signature device: %v", err)
		WriteAPIResponse(response, http.StatusNotFound, []string{
			err.Error(),
		})
		return
	}

	signature, err := signatureDevice.Sign(signDataRequest.Data)
	if err != nil {
		log.Printf("Error while signing data: %v", err)
		WriteInternalError(response)
		return
	}

	signDataResponse := SignDataResponse{
		Signature:  signature.Signature,
		SignedData: signature.Signed_Data,
	}
	WriteAPIResponse(response, http.StatusOK, signDataResponse)
}
