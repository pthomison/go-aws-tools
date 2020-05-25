package pkg

import (
	"errors"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
)

func (c *Client) FindInstancePrivateIP(instanceId string) (string, error) {
	inst, err := c.FindInstance(instanceId)

	if err != nil {
		return "", err
	}

	ptr := inst.PrivateIpAddress
	return *ptr, nil
}

func (c *Client) FindInstancePublicIP(instanceId string) (string, error) {
	inst, err := c.FindInstance(instanceId)

	if err != nil {
		return "", err
	}

	ptr := inst.PublicIpAddress
	return *ptr, nil
}

func (c *Client) FindInstanceAZ(instanceId string) (string, error) {
	inst, err := c.FindInstance(instanceId)

	if err != nil {
		return "", err
	}

	ptr := inst.Placement.AvailabilityZone
	return *ptr, nil
}

func (c *Client) FindInstance(instanceId string) (*ec2.Instance, error) {
	ec2_svc := ec2.New(c.sess)

	output_inst, err := ec2_svc.DescribeInstances(&ec2.DescribeInstancesInput{
		InstanceIds: []*string{
			aws.String(instanceId),
		},
	})

	if err != nil {
		handleAWSError(err)
		return nil, err
	}

	instance := output_inst.Reservations[0].Instances[0]

	return instance, nil
}

func (c *Client) FindInstanceIDByName(instanceName string) (string, error) {
	ec2_svc := ec2.New(c.sess)

	output_inst, err := ec2_svc.DescribeInstances(&ec2.DescribeInstancesInput{
		Filters: []*ec2.Filter{
			{
				Name:   aws.String("tag:Name"),
				Values: []*string{aws.String(instanceName)},
			},
		},
	})

	if err != nil {
		handleAWSError(err)
		return "", err
	}

	instance_count := len(output_inst.Reservations)

	if instance_count == 0 {
		return "", errors.New("No Instances Matched")
	} else if instance_count > 1 {
		return "", errors.New("Multiple Instances Matched")
	}

	return *output_inst.Reservations[0].Instances[0].InstanceId, nil
}
