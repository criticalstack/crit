package aws

import (
	"context"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/pkg/errors"
)

func GetInstanceInfo(ctx context.Context, cfg *aws.Config, instanceID string) (string, string, error) {
	sess, err := session.NewSession(cfg)
	if err != nil {
		return "", "", err
	}
	svc := ec2.New(sess)
	resp, err := svc.DescribeInstancesWithContext(ctx, &ec2.DescribeInstancesInput{
		InstanceIds: aws.StringSlice([]string{instanceID}),
	})
	if err != nil {
		return "", "", err
	}
	for _, r := range resp.Reservations {
		for _, instance := range r.Instances {
			ip := aws.StringValue(instance.PrivateIpAddress)
			profile := aws.StringValue(instance.IamInstanceProfile.Arn)
			return ip, profile, nil
		}
	}
	return "", "", errors.Errorf("instance not found: %v", instanceID)
}
