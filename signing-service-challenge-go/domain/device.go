package domain

import (
	"encoding/base64"
	"strconv"
	"sync"

	"github.com/fiskaly/coding-challenges/signing-service-challenge/crypto"
)

type Signature struct {
	Signature   string
	Signed_Data string
}

type SignatureDevice struct {
	Id                string
	Label             string
	signer            crypto.Signer
	signature_counter int
	last_signature    string
	mu                sync.Mutex
}

func NewSignatureDevice(id string, label string, signer crypto.Signer) *SignatureDevice {
	return &SignatureDevice{
		Id:     id,
		Label:  label,
		signer: signer,
	}
}

func (d *SignatureDevice) Sign(dataToBeSigned string) (*Signature, error) {
	d.mu.Lock()
	defer d.mu.Unlock()
	secured_data := d.getSecuredData(dataToBeSigned)
	signature, err := d.signer.Sign(secured_data)
	if err != nil {
		return nil, err
	}
	d.signature_counter++
	d.last_signature = base64.StdEncoding.EncodeToString(signature)

	return &Signature{
		Signature:   d.last_signature,
		Signed_Data: string(secured_data),
	}, nil
}

func (d *SignatureDevice) getSecuredData(dataToBeSigned string) []byte {
	var last_signature string

	if d.signature_counter == 0 {
		last_signature = base64.StdEncoding.EncodeToString([]byte(d.Id))
	} else {
		last_signature = d.last_signature
	}

	return []byte(strconv.Itoa(d.signature_counter) + "_" + dataToBeSigned + "_" + last_signature)
}
