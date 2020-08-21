package pki

import (
	"crypto"
	"crypto/x509"
	"path/filepath"

	"github.com/pkg/errors"
)

type KeyPair struct {
	Name string
	Cert *x509.Certificate
	Key  crypto.Signer
}

func (k *KeyPair) WriteFiles(dir string) error {
	return WriteCertAndKey(dir, k.Name, k.Cert, k.Key)
}

func LoadKeyPair(path, name string) (*KeyPair, error) {
	cert, err := ReadCertFromFile(filepath.Join(path, name+".crt"))
	if err != nil {
		return nil, err
	}
	key, err := ReadKeyFromFile(filepath.Join(path, name+".key"))
	if err != nil {
		return nil, err
	}
	kp := &KeyPair{
		Name: name,
		Cert: cert,
		Key:  key,
	}
	return kp, nil
}

type CertificateAuthority struct {
	*KeyPair
}

func NewCertificateAuthority(name string, cfg *Config) (*CertificateAuthority, error) {
	key, err := NewPrivateKey()
	if err != nil {
		return nil, errors.Wrap(err, "unable to create private key while generating CA certificate")
	}
	cert, err := NewSelfSignedCACert(cfg, key)
	if err != nil {
		return nil, errors.Wrap(err, "unable to create self-signed CA certificate")
	}
	ca := &CertificateAuthority{
		KeyPair: &KeyPair{
			Name: name,
			Cert: cert,
			Key:  key,
		},
	}
	return ca, nil
}

// NewSignedKeyPair returns a new KeyPair signed by the CA.
func (c *CertificateAuthority) NewSignedKeyPair(name string, cfg *Config) (*KeyPair, error) {
	key, err := NewPrivateKey()
	if err != nil {
		return nil, errors.Wrap(err, "unable to create private key")
	}
	cert, err := NewSignedCert(cfg, key, c.Cert, c.Key)
	if err != nil {
		return nil, errors.Wrap(err, "unable to sign certificate")
	}
	return &KeyPair{Name: name, Cert: cert, Key: key}, nil
}

func LoadCertificateAuthority(path, name string) (*CertificateAuthority, error) {
	kp, err := LoadKeyPair(path, name)
	if err != nil {
		return nil, err
	}
	return &CertificateAuthority{KeyPair: kp}, nil
}
