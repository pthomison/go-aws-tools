# go-aws-tools
This is repo is the result of a quick and dirty exploration golang, the AWS sdk, & cobra
I had implemented most of this functionality in bash, but turns out dealing with complicated data structs in a string based interpreter is a bad idea (shocking I know)

## Commands

A true README/man is tbd, but luckily cobra has a decent help menu built in. For the time being, here's the quick version 

	Usage:
	  go-aws-tools [command]

	Available Commands:
	  help           Help about any command
	  jump           Utility to authenticate & connect to selected instances w/ an in-memory key && ec2 instance connect
	  list-instances A very untested command to list instances w/ id, name, pub ip, priv ip

	Flags:
	  -h, --help             help for go-aws-tools
	      --profile string   aws profile to use
	      --region string    aws region to use (default "us-west-2")

	Use "go-aws-tools [command] --help" for more information about a command.


**jump**: authenticate & ssh into an instance; can tunnel through a bastion

	Utility to authenticate & connect to selected instances w/ an in-memory key && ec2 instance connect

	Usage:
	  go-aws-tools jump [flags]

	Flags:
	      --bastion-name string   if present, will attempt to tunnel through bastion
	  -h, --help                  help for jump
	      --id string             instance selector; mutually exclusive with name
	      --name string           instance selector; mutually exclusive with id
	      --user string           user override (default "ec2-user")

	Global Flags:
	      --profile string   aws profile to use
	      --region string    aws region to use (default "us-west-2")


**list-instances**: list out all running/pending instances in an account. Prints Name, ID, public IP, private IP; No Args
	A very untested command to list instances w/ id, name, pub ip, priv ip

	Usage:
	  go-aws-tools list-instances [flags]

	Flags:
	  -h, --help   help for list-instances

	Global Flags:
	      --profile string   aws profile to use
	      --region string    aws region to use (default "us-west-2")

## Building && Development

Refer to the makefile for a list of commands. If you don't have a working go environment, just use the docker wrappers

