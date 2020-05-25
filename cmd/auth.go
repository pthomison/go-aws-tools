package cmd

import (
	awsUtils "github.com/pthomison/go-aws-tools/pkg"
	"github.com/spf13/cobra"
	flag "github.com/spf13/pflag"
)

const (
	authPubKeyFlag  = "pubkey"
	authPrivKeyFlag = "privkey"
	authNameFlag    = "name"
	authIdFlag      = "id"
	authUserFlag    = "user"
)

var authCmd = &cobra.Command{
	Use:   "auth",
	Short: "Authenticates to an ec2-instance over instance-connect",
	Long:  ``,
	Run:   authCobra,
	Args:  cobra.ExactArgs(0),
}

func init() {
	// Attach Subcommand
	rootCmd.AddCommand(authCmd)

	// Create Subcommand Flags
	authCmd.PersistentFlags().String(authPubKeyFlag, "", "")
	authCmd.PersistentFlags().String(authPrivKeyFlag, "", "")
	authCmd.PersistentFlags().String(authNameFlag, "", "")
	authCmd.PersistentFlags().String(authIdFlag, "", "")
	authCmd.PersistentFlags().String(authUserFlag, "", "")
}

func authCobra(cmd *cobra.Command, args []string) {
	// Resolve Flags
	pubKeyF := cmd.Flags().Lookup(authPubKeyFlag)
	privKeyF := cmd.Flags().Lookup(authPrivKeyFlag)
	nameF := cmd.Flags().Lookup(authNameFlag)
	idF := cmd.Flags().Lookup(authIdFlag)
	userF := cmd.Flags().Lookup(authUserFlag)

	// mutually exclussive flag checking
	mutualExclusiveFlag(cmd, nameF, idF)
	mutualExclusiveFlag(cmd, pubKeyF, privKeyF)

	// Needed for if block initialization
	// var  string
	var err error

	// initialize client
	client, err := awsUtils.InitializeClient(awsProfile, awsRegion)
	if err != nil {
		handleGenericError(err)
		return
	}

	// determine instance id
	instanceId, err := resolveInstanceName(client, nameF, idF)
	if err != nil {
		handleGenericError(err)
		return
	}

	// _, rsaPublicKey, err := awsUtils.GenerateInMemoryKey()

	rsaPublicKey, err := loadPubKey(pubKeyF.Value.String())

	if err != nil {
		handleGenericError(err)
		return
	}

	client.Authenticate(instanceId, rsaPublicKey, userF.Value.String())
}
