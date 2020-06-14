package client

import (
	"crypto/rand"
	"crypto/rsa"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2instanceconnect"

	"golang.org/x/crypto/ssh"
)

// Sends Specified Key To instanceID; Handles AZ lookup
func (c *Client) Authenticate(instanceId string, pubKey *ssh.PublicKey, user string) error {
	// initialize the service
	svc := ec2instanceconnect.New(c.sess)

	// lookup AZ
	az, err := c.FindInstanceAZ(instanceId)
	if err != nil {
		return err
	}

	// stringify the pubkey
	pub := string(ssh.MarshalAuthorizedKey(*pubKey))

	// send key
	_, err = svc.SendSSHPublicKey(&ec2instanceconnect.SendSSHPublicKeyInput{
		AvailabilityZone: aws.String(az),
		InstanceId:       aws.String(instanceId),
		InstanceOSUser:   aws.String(user),
		SSHPublicKey:     aws.String(pub),
	})

	return err
}

// Creates in memory key for disposable authentication
func GenerateInMemoryKey(size int) (*rsa.PrivateKey, *ssh.PublicKey, error) {
	// generate key
	privateKey, err := rsa.GenerateKey(rand.Reader, size)
	if err != nil {
		return nil, nil, err
	}

	// strip public key && generate an ssh pub key
	rsaPublicKey := privateKey.Public()
	pubKey, err := ssh.NewPublicKey(rsaPublicKey)
	if err != nil {
		return nil, nil, err
	}

	return privateKey, &pubKey, nil
}
