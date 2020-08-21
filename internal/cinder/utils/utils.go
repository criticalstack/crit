package utils

import (
	"bytes"
	"compress/gzip"
	"encoding/pem"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/kyokomi/emoji"
	"k8s.io/client-go/util/keyutil"

	"github.com/criticalstack/crit/pkg/kubernetes/pki"
)

func CreateCA(cn string) ([]byte, []byte, error) {
	ca, err := pki.NewCertificateAuthority(cn, &pki.Config{
		CommonName: cn,
	})
	if err != nil {
		return nil, nil, err
	}
	encodedKey, err := keyutil.MarshalPrivateKeyToPEM(ca.Key)
	if err != nil {
		return nil, nil, err
	}
	block := pem.Block{
		Type:  "CERTIFICATE",
		Bytes: ca.Cert.Raw,
	}
	return pem.EncodeToMemory(&block), encodedKey, nil
}

func NewStep(msg string, icon string, verbose bool, fn func() error) error {
	suffix := fmt.Sprintf("  %s", msg)
	s := NewSpinner(os.Stdout)
	s.SetPrefix(" ")
	s.SetSuffix(suffix)
	if !verbose {
		s.Start()
	}
	defer s.Stop()
	if err := fn(); err != nil {
		if !verbose {
			_, _ = s.Write([]byte(emoji.Sprintf(" :cross_mark: %s\n", msg)))
		}
		return err
	}
	if !verbose {
		_, _ = s.Write([]byte(emoji.Sprintf(" %s %s\n", icon, msg)))
	}
	return nil
}

func Gunzip(data []byte) ([]byte, error) {
	r, err := gzip.NewReader(bytes.NewReader(data))
	if err != nil {
		return nil, err
	}
	defer func() { _ = r.Close() }()
	return ioutil.ReadAll(r)
}
