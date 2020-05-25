package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"

	awsUtils "github.com/pthomison/go-aws-tools/pkg"
)

// cobra doesn't appear to support this ootb, so this is a quick and dirty check w/ error
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
