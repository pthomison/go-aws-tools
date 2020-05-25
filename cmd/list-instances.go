package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	// "github.com/aws/aws-sdk-go/aws/session"
	// "github.com/aws/aws-sdk-go/service/ec2"
	// "github.com/aws/aws-sdk-go/aws/awserr"
	// "github.com/aws/aws-sdk-go/aws"

	awsUtils "github.com/pthomison/go-aws-tools/pkg"
)

// listInstancesCmd represents the listInstances command
var listInstancesCmd = &cobra.Command{
	Use:   "list-instances",
	Short: "A brief description of your command",
	Long: `A longer description
	test one
	test two`,
	Run: listInstances,
}

func init() {
	rootCmd.AddCommand(listInstancesCmd)
}

func listInstances(cmd *cobra.Command, args []string) {
	fmt.Println("list-instances called")

	// initialize client
	client, err := awsUtils.InitializeClient(awsProfile, awsRegion)
	if err != nil {
		handleGenericError(err)
		return
	}

	result, err := client.ListInstances()
	if err != nil {
		handleGenericError(err)
		return
	}

	for _, v := range result {
		fmt.Println(v.Str())
	}
}
