package pki

import (
	"crypto/sha256"
	"crypto/x509"
	"encoding/pem"
	"io/ioutil"

	"github.com/pkg/errors"
)

func GenerateCertHash(data []byte) ([]byte, error) {
	block, _ := pem.Decode(data)
	if block == nil {
		return nil, errors.New("cannot parse PEM formatted block")
	}
	cert, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		return nil, err
	}
	h := sha256.Sum256(cert.RawSubjectPublicKeyInfo)
	return h[:], nil
}

func GenerateCertHashFromFile(caCertPath string) ([]byte, error) {
	data, err := ioutil.ReadFile(caCertPath)
	if err != nil {
		return nil, err
	}
	return GenerateCertHash(data)
}
