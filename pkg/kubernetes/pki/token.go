package pki

import (
	"crypto/rand"
)

const validBootstrapTokenChars = "0123456789abcdefghijklmnopqrstuvwxyz"

// GenerateBootstrapToken constructs a bootstrap token in conformance with the
// following format:
// https://kubernetes.io/docs/admin/bootstrap-tokens/#token-format
func GenerateBootstrapToken() (id string, secret string) {
	token := make([]byte, 6+16)
	if _, err := rand.Read(token); err != nil {
		panic(err)
	}
	for i, b := range token {
		token[i] = validBootstrapTokenChars[int(b)%len(validBootstrapTokenChars)]
	}
	return string(token[:6]), string(token[6:])
}
