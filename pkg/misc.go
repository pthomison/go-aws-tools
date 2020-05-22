package pkg

import(
	"fmt"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/pthomison/go-aws-tools/internal"

	"github.com/aws/aws-sdk-go/aws/awserr"
	"errors"

)

type client struct {
	sess *session.Session
}

func (c *client) FindInstancePrivateIP(instanceId string) (string, error) {
	inst, _ := c.FindInstance(instanceId)
	ptr := inst.PrivateIpAddress
	return *ptr, nil
}

func (c *client) FindInstancePublicIP(instanceId string) (string, error) {
	inst, _ := c.FindInstance(instanceId)
	ptr := inst.PublicIpAddress
	return *ptr, nil
}

func (c *client) FindInstanceAZ(instanceId string) (string, error) {
	inst, _ := c.FindInstance(instanceId)
	ptr := inst.Placement.AvailabilityZone
	return *ptr, nil
}

func (c *client) FindInstance(instanceId string) (*ec2.Instance, error) {
	ec2_svc := ec2.New(c.sess)

	output_inst, err := ec2_svc.DescribeInstances( &ec2.DescribeInstancesInput{
	    InstanceIds: []*string{
	        aws.String(instanceId),
	    },
	})

	if err != nil {
	    internal.HandleAWSError(err)
	    return nil, err
	}

	instance := output_inst.Reservations[0].Instances[0]
	
	return instance, nil
}

func (c *client) FindInstanceByName(instanceName string) (string, error) {
	ec2_svc := ec2.New(c.sess)

	output_inst, err := ec2_svc.DescribeInstances( &ec2.DescribeInstancesInput{
	    Filters: []*ec2.Filter{
	        {
	            Name:   aws.String("tag:Name"),
	            Values: []*string{aws.String(instanceName)},
	        },
	    },
	})

	if err != nil {
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


func InitializeClient(awsProfile string, awsRegion string) *client {
	var c client
	s, _ := session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
		Profile: awsProfile,
		Config: aws.Config{
		    Region: aws.String(awsRegion),
		},
  	})
  	c.sess = s
  	return &c
}