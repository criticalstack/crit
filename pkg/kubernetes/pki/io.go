package pki

import (
	"crypto"
	"crypto/ecdsa"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"path/filepath"
	"time"

	"github.com/pkg/errors"
	certutil "k8s.io/client-go/util/cert"
	"k8s.io/client-go/util/keyutil"
)

func ReadCertFromFile(path string) (*x509.Certificate, error) {
	certs, err := certutil.CertsFromFile(path)
	if err != nil {
		return nil, err
	}

	// We are only putting one certificate in the certificate pem file, so it's safe to just pick the first one
	// TODO: Support multiple certs here in order to be able to rotate certs
	caCert := certs[0]

	// Check so that the certificate is valid now
	now := time.Now()
	if now.Before(caCert.NotBefore) {
		return nil, errors.New("the certificate is not valid yet")
	}
	if now.After(caCert.NotAfter) {
		return nil, errors.New("the certificate has expired")
	}
	return caCert, nil
}

func ReadKeyFromFile(path string) (crypto.Signer, error) {
	// Parse the private key from a file
	caKey, err := keyutil.PrivateKeyFromFile(path)
	if err != nil {
		return nil, errors.Wrapf(err, "couldn't load the private key file %s", path)
	}

	// Allow RSA and ECDSA formats only
	var key crypto.Signer
	switch k := caKey.(type) {
	case *rsa.PrivateKey:
		key = k
	case *ecdsa.PrivateKey:
		key = k
	default:
		return nil, errors.Errorf("the private key file %s is neither in RSA nor ECDSA format", path)
	}
	return key, nil
}

func WriteCertAndKey(path, name string, cert *x509.Certificate, key crypto.Signer) error {
	if err := WriteKey(path, name, key); err != nil {
		return err
	}
	return WriteCert(path, name, cert)

}

func WriteKey(path, name string, key crypto.Signer) error {
	encoded, err := keyutil.MarshalPrivateKeyToPEM(key)
	if err != nil {
		return err
	}
	return keyutil.WriteKey(filepath.Join(path, name+".key"), encoded)
}

// EncodeCertPEM returns PEM-endcoded certificate data
func EncodeCertPEM(cert *x509.Certificate) []byte {
	block := pem.Block{
		Type:  "CERTIFICATE",
		Bytes: cert.Raw,
	}
	return pem.EncodeToMemory(&block)
}

func MustEncodePrivateKeyPem(key crypto.Signer) []byte {
	data, err := keyutil.MarshalPrivateKeyToPEM(key)
	if err != nil {
		panic(err)
	}
	return data
}

func WriteCert(path, name string, cert *x509.Certificate) error {
	data := pem.EncodeToMemory(&pem.Block{
		Type:  "CERTIFICATE",
		Bytes: cert.Raw,
	})
	return certutil.WriteCert(filepath.Join(path, name+".crt"), data)
}

func WritePublicKey(path, name string, key crypto.PublicKey) error {
	der, err := x509.MarshalPKIXPublicKey(key)
	if err != nil {
		return err
	}
	data := pem.EncodeToMemory(&pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: der,
	})
	return keyutil.WriteKey(filepath.Join(path, name+".pub"), data)
}
