package cmd

// WIP

// const (
// 	publicKeyHeader = "ssh-rsa"
// )

// import (
// 	awsUtils "github.com/pthomison/go-aws-tools/pkg"
// 	"github.com/spf13/cobra"
// 	flag "github.com/spf13/pflag"
// )

// const (
// 	authPubKeyFlag  = "pubkey"
// 	authPrivKeyFlag = "privkey"
// 	authNameFlag    = "name"
// 	authIdFlag      = "id"
// 	authUserFlag    = "user"
// )

// var authCmd = &cobra.Command{
// 	Use:   "auth",
// 	Short: "Authenticates to an ec2-instance over instance-connect",
// 	Long:  ``,
// 	Run:   authCobra,
// 	Args:  cobra.ExactArgs(0),
// }

// func init() {
// 	// Attach Subcommand
// 	rootCmd.AddCommand(authCmd)

// 	// Create Subcommand Flags
// 	authCmd.PersistentFlags().String(authPubKeyFlag, "", "")
// 	authCmd.PersistentFlags().String(authPrivKeyFlag, "", "")
// 	authCmd.PersistentFlags().String(authNameFlag, "", "")
// 	authCmd.PersistentFlags().String(authIdFlag, "", "")
// 	authCmd.PersistentFlags().String(authUserFlag, "", "")
// }

// func authCobra(cmd *cobra.Command, args []string) {
// 	// Resolve Flags
// 	pubKeyF := cmd.Flags().Lookup(authPubKeyFlag)
// 	privKeyF := cmd.Flags().Lookup(authPrivKeyFlag)
// 	nameF := cmd.Flags().Lookup(authNameFlag)
// 	idF := cmd.Flags().Lookup(authIdFlag)
// 	userF := cmd.Flags().Lookup(authUserFlag)

// 	// mutually exclussive flag checking
// 	mutualExclusiveFlag(cmd, nameF, idF)
// 	mutualExclusiveFlag(cmd, pubKeyF, privKeyF)

// 	// Needed for if block initialization
// 	// var  string
// 	var err error

// 	// initialize client
// 	client, err := awsUtils.InitializeClient(awsProfile, awsRegion)
// 	if err != nil {
// 		handleGenericError(err)
// 		return
// 	}

// 	// determine instance id
// 	instanceId, err := resolveInstanceName(client, nameF, idF)
// 	if err != nil {
// 		handleGenericError(err)
// 		return
// 	}

// 	// _, rsaPublicKey, err := awsUtils.GenerateInMemoryKey()

// 	rsaPublicKey, err := loadPubKey(pubKeyF.Value.String())

// 	if err != nil {
// 		handleGenericError(err)
// 		return
// 	}

// 	client.Authenticate(instanceId, rsaPublicKey, userF.Value.String())
// }

// func loadPubKey(filename string) (*ssh.PublicKey, error) {
// 	fmt.Printf("%+v\n", filename)

// 	pub, err := ioutil.ReadFile(filename)
// 	if err != nil {
// 		return nil, err
// 	}

// 	s := strings.Split(string(pub), " ")
// 	pubKeyData := s[1]

// 	byteArray := []byte(pubKeyData)
// 	byteDecoded := make([]byte, len(byteArray))

// 	if _, err := base64.StdEncoding.Decode(byteDecoded, byteArray); err != nil {
// 		return nil, err
// 	}

// 	lenLen := 4
// 	keyTypeLenStart := 0

// 	keyTypeStart := keyTypeLenStart + lenLen

// 	keyTypeLen := int(binary.BigEndian.Uint32(byteDecoded[keyTypeLenStart:keyTypeStart]))

// 	eLenStart := keyTypeStart + keyTypeLen
// 	eStart := eLenStart + lenLen

// 	eLen := int(binary.BigEndian.Uint32(byteDecoded[eLenStart:eStart]))

// 	nLenStart := eStart + eLen
// 	nStart := nLenStart + lenLen

// 	nLen := int(binary.BigEndian.Uint32(byteDecoded[nLenStart:nStart]))

// 	keyType := byteDecoded[keyTypeStart:eLenStart]
// 	e := byteDecoded[eStart:nLenStart]
// 	n := byteDecoded[nStart:]

// 	var bigN, bigE big.Int

// 	eNum, eNumErr := binary.Uvarint(e)

// 	pubKey, err := ssh.NewPublicKey(&rsa.PublicKey{
// 		N: bigN.SetBytes(n),
// 		E: int(bigE.SetBytes(e).Uint64()),
// 	})

// 	_ = keyType
// 	_ = nLen
// 	_ = eNumErr
// 	_ = eNum

// 	return &pubKey, err
// }

// func loadPrivKey(filename string) *rsa.PrivateKey {
// 	return nil
// }
