package cmd

import (
	"fmt"
	"io/ioutil"
	"os"
	// "encoding/pem"
	// "crypto/x509"
	"crypto/rsa"
	"encoding/base64"
	"github.com/spf13/pflag"
	"golang.org/x/crypto/ssh"
	"strings"

	"encoding/binary"

	"math/big"
	// "bufio"
	// "bytes"
	"github.com/spf13/cobra"

	awsUtils "github.com/pthomison/go-aws-tools/pkg"
)

const (
	publicKeyHeader = "ssh-rsa"
)

func handleGenericError(err error) {
	fmt.Println(err)
	return
}

func mutualExclusiveFlag(cmd *cobra.Command, flagA *pflag.Flag, flagB *pflag.Flag) error {
	if !flagA.Changed && !flagB.Changed {
		cmd.Help()
		return fmt.Errorf("Must supply a target")
	} else if flagA.Changed && flagB.Changed {
		cmd.Help()
		return fmt.Errorf("Must supply a single target")
	}
	return nil
}

func loadPubKey(filename string) (*ssh.PublicKey, error) {
	fmt.Printf("%+v\n", filename)

	pub, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	s := strings.Split(string(pub), " ")
	pubKeyData := s[1]

	byteArray := []byte(pubKeyData)
	byteDecoded := make([]byte, len(byteArray))

	if _, err := base64.StdEncoding.Decode(byteDecoded, byteArray); err != nil {
		return nil, err
	}

	lenLen := 4
	keyTypeLenStart := 0

	keyTypeStart := keyTypeLenStart + lenLen

	keyTypeLen := int(binary.BigEndian.Uint32(byteDecoded[keyTypeLenStart:keyTypeStart]))

	eLenStart := keyTypeStart + keyTypeLen
	eStart := eLenStart + lenLen

	eLen := int(binary.BigEndian.Uint32(byteDecoded[eLenStart:eStart]))

	nLenStart := eStart + eLen
	nStart := nLenStart + lenLen

	nLen := int(binary.BigEndian.Uint32(byteDecoded[nLenStart:nStart]))

	keyType := byteDecoded[keyTypeStart:eLenStart]
	e := byteDecoded[eStart:nLenStart]
	n := byteDecoded[nStart:]

	var bigN, bigE big.Int

	eNum, eNumErr := binary.Uvarint(e)

	pubKey, err := ssh.NewPublicKey(&rsa.PublicKey{
		N: bigN.SetBytes(n),
		E: int(bigE.SetBytes(e).Uint64()),
	})

	_ = keyType
	_ = nLen
	_ = eNumErr
	_ = eNum

	return &pubKey, err
}

func loadPrivKey(filename string) *rsa.PrivateKey {
	return nil
}

func commandError(e error) {
	if e != nil {
		fmt.Println(e)
		os.Exit(1)
	}
}

func resolveInstanceName(c *awsUtils.Client, nameFlag *flag.Flag, idFlag *flag.Flag) (string, error) {
	var instanceId string
	var err error
	if nameFlag.Changed {
		instanceId, err = c.FindInstanceByName(nameFlag.Value.String())
		if err != nil {
			return "", err
		}
	} else {
		instanceId = idFlag.Value.String()
	}
	return instanceId, nil
}
