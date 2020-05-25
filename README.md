# go-aws-tools
This is repo is the result of a quick and dirty exploration golang, the AWS sdk, & cobra
I had implemented most of this functionality in bash, but turns out dealing with complicated data structs in a string based interpreter is a bad idea (shocking I know)

## Commands

A true README/man is tbd, but luckily cobra has a decent help menu built in. For the time being, here's the quick version 

**jump**: authenticate & ssh into an instance; can tunnel through a bastion
**list-instances**: list out all running/pending instances in an account. Prints Name, ID, public IP, private IP