package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	awsUtils "github.com/pthomison/go-aws-tools/pkg/client"
)

// listInstancesCmd represents the listInstances command
var listInstancesCmd = &cobra.Command{
	Use:   "list-instances",
	Short: "A very untested command to list instances w/ id, name, pub ip, priv ip",
	Long:  ``,
	Run:   listInstances,
	Args:  cobra.ExactArgs(0),
}

func init() {
	rootCmd.AddCommand(listInstancesCmd)
}

func listInstances(cmd *cobra.Command, args []string) {
	fmt.Println("list-instances called")

	// initialize client
	client, err := awsUtils.InitializeClient(awsProfile, awsRegion)
	commandError(err)

	// Collect Instances (NOTE: using type from pkg, not ec2.Instance)
	result, err := client.ListInstances()
	commandError(err)

	// print each instance on a newline
	for _, v := range result {
		fmt.Println(v.Str())
	}
}
