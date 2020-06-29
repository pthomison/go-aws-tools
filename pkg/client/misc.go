package client

import (
	"errors"
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/ec2"
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
			var id, name, privateIP, publicIP string
			id = *inst.InstanceId
			name, err := findTagValue(inst, "Name")
			if err != nil {
				return nil, err
			}

			if inst.PrivateIpAddress != nil {
				privateIP = *inst.PrivateIpAddress
			}

			if inst.PublicIpAddress != nil {
				publicIP = *inst.PublicIpAddress
			}

			instanceDescriptions = append(instanceDescriptions, &instanceDescription{
				instanceId:       id,
				instanceName:     name,
				privateIpAddress: privateIP,
				publicIpAddress:  publicIP,
			})
		}
	}

	if err != nil {
		handleAWSError(err)
		return nil, err
	}

	return instanceDescriptions, nil
}

func (c *Client) NameResource(resourceId string, name string) error {
	svc := ec2.New(c.sess)

	tagInput := &ec2.CreateTagsInput{
		Resources: []*string{
			aws.String(resourceId),
		},
		Tags: []*ec2.Tag{
			{
				Key:   aws.String("Name"),
				Value: aws.String(name),
			},
		},
	}

	_, err := svc.CreateTags(tagInput)
	return err
}

type instanceDescription struct {
	instanceId       string
	instanceName     string
	privateIpAddress string
	publicIpAddress  string
}

func (desc *instanceDescription) Str() string {
	if desc.publicIpAddress == "" {
		return desc.instanceName + " : " + desc.instanceId + " : " + desc.privateIpAddress
	} else {
		return desc.instanceName + " : " + desc.instanceId + " : " + desc.privateIpAddress + " : " + desc.publicIpAddress
	}
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

func handleAWSError(err error) {
	if aerr, ok := err.(awserr.Error); ok {
		switch aerr.Code() {
		default:
			fmt.Println(aerr.Error())
		}
	} else {
		// Print the error, cast err to awserr.Error to get the Code and
		// Message from an error.
		fmt.Println(err.Error())
	}
	return
}

func tagLookup(tags []*ec2.Tag, key string) *ec2.Tag {
	for _, v := range tags {
		if *v.Key == key {
			return v
		}
	}
	return nil
}
