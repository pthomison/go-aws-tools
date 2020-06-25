package cmd

import (
	"github.com/spf13/cobra"

	awsUtils "github.com/pthomison/go-aws-tools/pkg/client"
)

const (
	jumpBastionFlag = "bastion-name"
	jumpNameFlag    = "name"
	jumpIdFlag      = "id"
	jumpUserFlag    = "user"
)

var jumpCmd = &cobra.Command{
	Use:   "jump",
	Short: "Utility to authenticate & connect to selected instances w/ an in-memory key && ec2 instance connect",
	Long:  ``,
	Run:   jumpCobra,
	Args:  cobra.ExactArgs(0),
}

func init() {
	rootCmd.AddCommand(jumpCmd)
	jumpCmd.PersistentFlags().String(jumpBastionFlag, "", "if present, will attempt to tunnel through bastion")
	jumpCmd.PersistentFlags().String(jumpNameFlag, "", "instance selector; mutually exclusive with "+jumpIdFlag)
	jumpCmd.PersistentFlags().String(jumpIdFlag, "", "instance selector; mutually exclusive with "+jumpNameFlag)
	jumpCmd.PersistentFlags().String(jumpUserFlag, "ec2-user", "user override")

}

func jumpCobra(cmd *cobra.Command, args []string) {
	bastionF := cmd.Flags().Lookup(jumpBastionFlag)
	nameF := cmd.Flags().Lookup(jumpNameFlag)
	idF := cmd.Flags().Lookup(jumpIdFlag)
	userF := cmd.Flags().Lookup(jumpUserFlag)

	// flag checking
	mutualExclusiveFlag(cmd, nameF, idF)

	// initialize client
	client, err := awsUtils.InitializeClient(awsProfile, awsRegion)
	commandError(err)

	// determine instance id
	instanceId, err := resolveInstanceName(client, nameF, idF)
	commandError(err)

	// generate temporary in memory key
	privateKey, rsaPublicKey, err := awsUtils.GenerateInMemoryKey(4096)
	commandError(err)

	// find user
	user := userF.Value.String()

	// if bastion flag present, tunnel
	if bastionF.Changed {
		bastionId, err := client.FindInstanceIDByName(bastionF.Value.String())
		commandError(err)

		// authenicate to both instances before tunneling
		client.Authenticate(instanceId, rsaPublicKey, user)
		client.Authenticate(bastionId, rsaPublicKey, user)

		// tunnel
		err = client.JumpThroughBastion(instanceId, bastionId, privateKey, user)
		commandError(err)
	} else {
		// authenticate && jump
		client.Authenticate(instanceId, rsaPublicKey, user)
		client.Jump(instanceId, privateKey, user)
		commandError(err)
	}

}
