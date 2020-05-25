package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
)

var cfgFile string
var awsProfile string
var awsRegion string

var rootCmd = &cobra.Command{
	Use:   "go-aws-tools",
	Short: "",
	Long:  ``,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringVar(&awsProfile, "profile", "", "aws profile to use")
	rootCmd.PersistentFlags().StringVar(&awsRegion, "region", "us-west-2", "aws region to use")
}
