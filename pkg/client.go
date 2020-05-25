package pkg

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
)

// object to hold the AWS session for successive calls
type Client struct {
	sess *session.Session
}

// initializes session && return client
func InitializeClient(awsProfile string, awsRegion string) (*Client, error) {
	s, err := session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
		Profile:           awsProfile,
		Config: aws.Config{
			Region: aws.String(awsRegion),
		},
	})
	if err != nil {
		return nil, err
	}

	return &Client{
		sess: s,
	}, nil
}
