package pkg

import (
	"crypto/rand"
	"crypto/rsa"

	"github.com/aws/aws-sdk-go/service/ec2instanceconnect"
	"github.com/aws/aws-sdk-go/aws"

	// "github.com/pthomison/go-aws-tools/internal/helper"
	"golang.org/x/crypto/ssh"
)

func (c *client) Authenticate(instanceId string, pubKey *ssh.PublicKey, user string) error {
	svc := ec2instanceconnect.New(c.sess)
	az, _ := c.FindInstanceAZ(instanceId)
	pub := string(ssh.MarshalAuthorizedKey(*pubKey))

	input := &ec2instanceconnect.SendSSHPublicKeyInput{
		AvailabilityZone: aws.String(az),
		InstanceId:       aws.String(instanceId),
		InstanceOSUser:   aws.String(user),
		SSHPublicKey:     aws.String(pub),
	}

	_, err := svc.SendSSHPublicKey(input)
	return err
}

func GenerateInMemoryKey() (*rsa.PrivateKey, *ssh.PublicKey, error) {
	privateKey, err := rsa.GenerateKey(rand.Reader, 4096)
	if err != nil {
		return nil, nil, err
	}
	rsaPublicKey := privateKey.Public()
	pubKey, err := ssh.NewPublicKey(rsaPublicKey)
	if err != nil {
		return nil, nil, err
	}
	return privateKey, &pubKey, nil
}