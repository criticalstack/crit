package bootstrap

import "encoding/json"

type AuthorizationType int

const (
	UnknownAuthorizationType AuthorizationType = iota
	AmazonIdentityDocumentAndSignature
)

func (at AuthorizationType) String() string {
	switch at {
	case AmazonIdentityDocumentAndSignature:
		return "aws"
	default:
		return "unknown"
	}
}

func (at AuthorizationType) MarshalText() ([]byte, error) {
	switch at {
	case AmazonIdentityDocumentAndSignature:
		return []byte("aws"), nil
	default:
		return []byte("unknown"), nil
	}
}

func (at *AuthorizationType) UnmarshalText(data []byte) error {
	switch string(data) {
	case "aws":
		*at = AmazonIdentityDocumentAndSignature
	default:
		*at = UnknownAuthorizationType
	}
	return nil
}

type Request struct {
	Type AuthorizationType `json:"type"`
	Body json.RawMessage   `json:"body"`
}

type Response struct {
	Error          string `json:"error"`
	BootstrapToken string `json:"bootstrapToken"`
}
