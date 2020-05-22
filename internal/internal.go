package internal

import(
	"fmt"
	"github.com/aws/aws-sdk-go/aws/awserr"
)

func HandleGenericError(err error) {
	fmt.Println(err)
	return
}

func HandleAWSError(err error) {
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