package cmd

import(
	"fmt"
	"io/ioutil"
	// "os"
	// "encoding/pem"
	// "crypto/x509"
	// "crypto/rsa"
	"github.com/spf13/pflag"
	"golang.org/x/crypto/ssh"
	"encoding/base64"
	"strings"

	"encoding/binary"

)

const (
	publicKeyHeader = "ssh-rsa"
)



func handleGenericError(err error) {
	fmt.Println(err)
	return
}

func mutualExclusiveFlag(flagA *pflag.Flag, flagB *pflag.Flag) error {
	if !flagA.Changed && !flagB.Changed {
		return fmt.Errorf("Must supply a target") 
	} else if flagA.Changed && flagB.Changed {
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

	byteArray :=  []byte(pubKeyData)
	byteDecoded := make([]byte, len(byteArray))

	if _, err := base64.StdEncoding.Decode(byteDecoded, byteArray); err != nil {
		return nil, err
	}

	lenLen := 4
	keyTypeLenStart := 0

	keyTypeStart    := keyTypeLenStart + lenLen

	keyTypeLen      := int(binary.BigEndian.Uint32(byteDecoded[keyTypeLenStart:keyTypeStart]))

	eLenStart       := keyTypeStart + keyTypeLen
	eStart          := eLenStart + lenLen

	eLen            := int(binary.BigEndian.Uint32(byteDecoded[eLenStart:eStart]))

	nLenStart       := eStart + eLen
	nStart          := nLenStart + lenLen

	nLen            := int(binary.BigEndian.Uint32(byteDecoded[nLenStart:nStart]))

	keyType         := byteDecoded[keyTypeStart:eLenStart]
	e               := byteDecoded[eStart:nLenStart]
	n               := byteDecoded[nStart:]


	fmt.Printf("key type len: %+v\n", keyTypeLen)
	fmt.Printf("key type: %+v\n", keyType)

	fmt.Printf("e len: %+v\n", keyTypeLen)
	fmt.Printf("e: %+v\n", e)

	fmt.Printf("n len: %+v\n", eLen)
	fmt.Printf("n: %+v\n", n)

	// for i, v := range byteDecoded {
	// 	for i:=0;i<len(publicKeyHeader)
	// }

	// encoding


	// pubBlock, _ := pem.Decode(pub)

	// fmt.Printf("%+v\n", pubBlock)

	// pubByte, err := x509.ParsePKIXPublicKey(pub)
	// if err != nil {
	// 	return nil, err
	// }

	// // pubKey, err := ssh.ParsePublicKey([]byte(pubKeyData))
	// var rsaKey rsa.PublicKey

	// fmt.Printf("%+v\n", ssh.Unmarshal(pub, &rsaKey))

	pubKey, err := ssh.ParsePublicKey(pub)

	return &pubKey, err
}