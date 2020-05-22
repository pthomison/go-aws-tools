package cmd

import (
	"fmt"
	"github.com/spf13/cobra"

	awsUtils "github.com/pthomison/go-aws-tools/pkg"
	"github.com/pthomison/go-aws-tools/internal"
)

const (
	bastionFlag = "bastion"
	nameFlag    = "name"
	idFlag      = "id"
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
	jumpCmd.PersistentFlags().String(bastionFlag, "", "")
	jumpCmd.PersistentFlags().String(nameFlag, "", "")
	jumpCmd.PersistentFlags().String(idFlag, "", "")
}

func jumpCobra(cmd *cobra.Command, args []string) {
	bastionF := cmd.Flags().Lookup(bastionFlag)
	nameF    := cmd.Flags().Lookup(nameFlag)
	idF      := cmd.Flags().Lookup(idFlag)

	// mutually exclussive flag checking
	if !nameF.Changed && !idF.Changed {
		internal.HandleGenericError(fmt.Errorf("Must supply a target"))
		return
	} else if nameF.Changed && idF.Changed {
		internal.HandleGenericError(fmt.Errorf("Must supply a single target"))
		return
	}

	var instanceId string
	var bastionId string
	var err error

	// initialize client
	client := awsUtils.InitializeClient(awsProfile, awsRegion)

	// determine instance id
	if nameF.Changed {
		instanceId, err = client.FindInstanceByName(nameF.Value.String())

		if err != nil {
			internal.HandleGenericError(err)
			return
		}
	} else {
		instanceId = idF.Value.String()
	}

	// generate temporary in memory key
	privateKey, rsaPublicKey, err := awsUtils.GenerateInMemoryKey()
	if err != nil {
		internal.HandleGenericError(err)
		return
	}

	if bastionF.Changed {
		bastionId, err = client.FindInstanceByName(bastionF.Value.String())
		if err != nil {
			internal.HandleGenericError(err)
			return
		}

		client.Authenticate(instanceId, rsaPublicKey, "ec2-user")
		client.Authenticate(bastionId, rsaPublicKey, "ec2-user")

		if err = client.JumpThroughBastion(instanceId, bastionId, privateKey, "ec2-user"); err != nil {
			internal.HandleGenericError(err)
			return
		}
	} else {
		client.Authenticate(instanceId, rsaPublicKey, "ec2-user")


		if err = client.Jump(instanceId, privateKey, "ec2-user"); err != nil {
			internal.HandleGenericError(err)
			return
		}
	}

}