package cmd

import (
	"github.com/spf13/cobra"

	awsUtils "github.com/pthomison/go-aws-tools/pkg"
)

const (
	jumpBastionFlag = "bastion"
	jumpNameFlag    = "name"
	jumpIdFlag      = "id"
)

var jumpCmd = &cobra.Command{
	Use:   "jump",
	Short: "",
	Long: ``,
	Run: jumpCobra,
	Args: cobra.ExactArgs(0),
}

func init() {
	rootCmd.AddCommand(jumpCmd)
	jumpCmd.PersistentFlags().String(jumpBastionFlag, "", "")
	jumpCmd.PersistentFlags().String(jumpNameFlag, "", "")
	jumpCmd.PersistentFlags().String(jumpIdFlag, "", "")
}

func jumpCobra(cmd *cobra.Command, args []string) {
	bastionF := cmd.Flags().Lookup(jumpBastionFlag)
	nameF    := cmd.Flags().Lookup(jumpNameFlag)
	idF      := cmd.Flags().Lookup(jumpIdFlag)

	// mutually exclussive flag checking
	mutualExclusiveFlag(nameF, idF)

	var instanceId string
	var bastionId string
	var err error

	// initialize client
	client := awsUtils.InitializeClient(awsProfile, awsRegion)

	// determine instance id
	if nameF.Changed {
		instanceId, err = client.FindInstanceByName(nameF.Value.String())

		if err != nil {
			handleGenericError(err)
			return
		}
	} else {
		instanceId = idF.Value.String()
	}

	// generate temporary in memory key
	privateKey, rsaPublicKey, err := awsUtils.GenerateInMemoryKey()
	if err != nil {
		handleGenericError(err)
		return
	}

	if bastionF.Changed {
		bastionId, err = client.FindInstanceByName(bastionF.Value.String())
		if err != nil {
			handleGenericError(err)
			return
		}

		client.Authenticate(instanceId, rsaPublicKey, "ec2-user")
		client.Authenticate(bastionId, rsaPublicKey, "ec2-user")

		if err = client.JumpThroughBastion(instanceId, bastionId, privateKey, "ec2-user"); err != nil {
			handleGenericError(err)
			return
		}
	} else {
		client.Authenticate(instanceId, rsaPublicKey, "ec2-user")


		if err = client.Jump(instanceId, privateKey, "ec2-user"); err != nil {
			handleGenericError(err)
			return
		}
	}

}