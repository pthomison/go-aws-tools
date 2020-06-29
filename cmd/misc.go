package cmd

import (
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"

	awsUtils "github.com/pthomison/go-aws-tools/pkg/client"
)

// cobra doesn't appear to support this ootb, so this is a quick and dirty check w/ error
func mutualExclusiveFlag(cmd *cobra.Command, flagA *pflag.Flag, flagB *pflag.Flag) {
	if !flagA.Changed && !flagB.Changed {
		cmd.Help()
		os.Exit(1)
	} else if flagA.Changed && flagB.Changed {
		cmd.Help()
		os.Exit(1)
	}
}

func requiredFlag(cmd *cobra.Command, flag *pflag.Flag) {
	if !flag.Changed {
		cmd.Help()
		os.Exit(1)
	}
}

// i got really tired of writing err!=nil blocks
func commandError(e error) {
	if e != nil {
		fmt.Printf("%+v\n", e)
		os.Exit(1)
	}
}

// useful block that gets called in a couple of the commands
// looks up name if provided; else returns id value
func resolveInstanceName(c *awsUtils.Client, nameFlag *pflag.Flag, idFlag *pflag.Flag) (string, error) {
	var instanceId string
	var err error
	if nameFlag.Changed {
		instanceId, err = c.FindInstanceIDByName(nameFlag.Value.String())
		if err != nil {
			return "", err
		}
	} else {
		instanceId = idFlag.Value.String()
	}
	return instanceId, nil
}

func tagLookup(tags []*ec2.Tag, key string) *ec2.Tag {
	for _, v := range tags {
		if *v.Key == key {
			return v
		}
	}
	return nil
}
