package cmd

import (
	"github.com/spf13/cobra"

	awsUtils "github.com/pthomison/go-aws-tools/pkg"
	"github.com/pthomison/go-aws-tools/internal"
)

var authCmd = &cobra.Command{
	Use:   "auth",
	Short: "",
	Long: ``,
	Run: authCobra,
	Args: cobra.ExactArgs(0),
}

func init() {
	rootCmd.AddCommand(authCmd)
}

func authCobra(cmd *cobra.Command, args []string) {
	instanceName := args[0]

	client := awsUtils.InitializeClient(awsProfile, awsRegion)
	instanceId, err := client.FindInstanceByName(instanceName)

	if err != nil {
		internal.HandleGenericError(err)
		return
	}

	_, rsaPublicKey, err := awsUtils.GenerateInMemoryKey()

	if err != nil {
		internal.HandleGenericError(err)
		return
	}

	client.Authenticate(instanceId, rsaPublicKey, "ec2-user")
}




