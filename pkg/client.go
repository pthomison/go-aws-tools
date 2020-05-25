package pkg

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
)

type Client struct {
	sess *session.Session
}

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
