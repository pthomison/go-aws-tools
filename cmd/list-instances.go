/*
Copyright © 2020 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws"
)

// listInstancesCmd represents the listInstances command
var listInstancesCmd = &cobra.Command{
	Use:   "list-instances",
	Short: "A brief description of your command",
	Long: `A longer description
	test one
	test two`,
	Run: listInstances,
}

func listInstances(cmd *cobra.Command, args []string) {
	fmt.Println("list-instances called")

		   //  InstanceIds: []*string{
	    //     aws.String("i-1234567890abcdef0"),
	    // },

	svc := ec2.New(session.New(), aws.NewConfig().WithRegion("us-west-2"))
	input := &ec2.DescribeInstancesInput{
	}

	result, err := svc.DescribeInstances(input)
	if err != nil {
	    if aerr, ok := err.(awserr.Error); ok {
	        switch aerr.Code() {
	        default:
	            fmt.Println(aerr.Error())
	        }
	    } else {
	        // Print the error, cast err to awserr.Error to get the Code and
	        // Message from an error.
	        fmt.Println(err.Error())
	    }
	    return
	}
	fmt.Println(result)
}

func init() {
	rootCmd.AddCommand(listInstancesCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// listInstancesCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// listInstancesCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
