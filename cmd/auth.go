package cmd

import (
	"github.com/spf13/cobra"

	awsUtils "github.com/pthomison/go-aws-tools/pkg"
)

const (
	authKeyFlag     = "pubkey"
	authNameFlag    = "name"
	authIdFlag      = "id"
	authUserFlag    = "user"
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
	authCmd.PersistentFlags().String(authKeyFlag, "", "")
	authCmd.PersistentFlags().String(authNameFlag, "", "")
	authCmd.PersistentFlags().String(authIdFlag, "", "")
	authCmd.PersistentFlags().String(authUserFlag, "", "")
}

func authCobra(cmd *cobra.Command, args []string) {
	keyF     := cmd.Flags().Lookup(authKeyFlag)
	nameF    := cmd.Flags().Lookup(authNameFlag)
	idF      := cmd.Flags().Lookup(authIdFlag)
	userF    := cmd.Flags().Lookup(authUserFlag)


	_, _ = keyF, userF

	// mutually exclussive flag checking
	mutualExclusiveFlag(nameF, idF)

	var instanceId string
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

	// _, rsaPublicKey, err := awsUtils.GenerateInMemoryKey()

	rsaPublicKey, err := loadPubKey(keyF.Value.String())

	if err != nil {
		handleGenericError(err)
		return
	}

	client.Authenticate(instanceId, rsaPublicKey, userF.Value.String())
}




