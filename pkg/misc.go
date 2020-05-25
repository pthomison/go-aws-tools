package pkg

import (
	"errors"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/pthomison/go-aws-tools/internal"
)

func (c *Client) ListInstances() ([]*instanceDescription, error) {
	svc := ec2.New(c.sess)

	instances, err := svc.DescribeInstances(&ec2.DescribeInstancesInput{
		Filters: []*ec2.Filter{
			{
				Name: aws.String("instance-state-name"),
				Values: []*string{
					aws.String("pending"),
					aws.String("running"),
				},
			},
		},
	})

	var instanceDescriptions []*instanceDescription

	for _, resv := range instances.Reservations {
		for _, inst := range resv.Instances {
			id := *inst.InstanceId
			name, err := findTagValue(inst, "Name")
			if err != nil {
				return nil, err
			}
			privateIP := *inst.PrivateIpAddress
			publicIP := *inst.PublicIpAddress
			instanceDescriptions = append(instanceDescriptions, &instanceDescription{
				instanceId:       id,
				instanceName:     name,
				privateIpAddress: privateIP,
				publicIpAddress:  publicIP,
			})
		}
	}

	if err != nil {
		internal.HandleAWSError(err)
		return nil, err
	}

	return instanceDescriptions, nil
}

type instanceDescription struct {
	instanceId       string
	instanceName     string
	privateIpAddress string
	publicIpAddress  string
}

func (desc *instanceDescription) Str() string {
	s := desc.instanceName + " : " + desc.instanceId + " : " + desc.privateIpAddress + " : " + desc.publicIpAddress
	return s
}

func findTagValue(instance *ec2.Instance, tagName string) (string, error) {
	tags := instance.Tags
	for _, v := range tags {
		if *v.Key == tagName {
			return *v.Value, nil
		}
	}

	return "", errors.New("Tag Not Found")
}
