package ec2metadata

import (
	"context"
	"crypto/sha1"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/ec2metadata"
	"github.com/aws/aws-sdk-go/aws/session"
)

// A Document provides a struct for EC2 instance identity documents to be
// unmarshaled.
type Document struct {
	AvailabilityZone string `json:"availabilityZone"`
	Architecture     string `json:"architecture"`
	PrivateIp        string `json:"privateIp"`
	Region           string `json:"region"`
	InstanceId       string `json:"instanceId"`
	AccountId        string `json:"accountId"`
	InstanceType     string `json:"instanceType"`
	ImageId          string `json:"imageId"`
}

// Name returns the logical name for the instance described in the identity
// document and is the value used when deriving the unique identifier hash.
func (d *Document) Name() string {
	return fmt.Sprintf("%s-%s", d.AccountId, d.InstanceId)
}

func (d *Document) Hash() string {
	return strings.ToUpper(fmt.Sprintf("%x", sha1.Sum([]byte(d.Name()))))
}

type SignedDocument struct {
	Document  []byte `json:"document"`
	Signature string `json:"signature"`
}

func GetSignedDocument() ([]byte, error) {
	sess, err := session.NewSession(aws.NewConfig())
	if err != nil {
		return nil, err
	}
	svc := ec2metadata.New(sess)
	doc, err := svc.GetDynamicDataWithContext(context.TODO(), "instance-identity/document")
	if err != nil {
		return nil, err
	}
	sig, err := svc.GetDynamicDataWithContext(context.TODO(), "instance-identity/signature")
	if err != nil {
		return nil, err
	}
	return json.Marshal(SignedDocument{
		Document:  []byte(doc),
		Signature: sig,
	})
}
