package client

import (
	"crypto/rsa"
	"fmt"
	"io"
	"os"

	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/terminal"
)

func (c *Client) Jump(instanceId string, privKey *rsa.PrivateKey, username string) error {
	signer, err := ssh.NewSignerFromKey(privKey)
	if err != nil {
		return fmt.Errorf("Can't Create Signer From Private Key: %w", err)
	}

	// Find IP of public instance
	publicIP, err := c.FindInstancePublicIP(instanceId)
	if err != nil {
		return fmt.Errorf("Can't Find Public IP on %v: %w", instanceId, err)
	}

	dialAddress := publicIP + ":22"

	client, err := ssh.Dial("tcp", dialAddress, &ssh.ClientConfig{
		User: username,
		Auth: []ssh.AuthMethod{
			ssh.PublicKeys(signer),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	})
	if err != nil {
		return fmt.Errorf("Failed to dial : %w", err)
	}

	return shell(client)
}

func (c *Client) JumpThroughBastion(instanceId string, bastionId string, privKey *rsa.PrivateKey, username string) error {
	signer, err := ssh.NewSignerFromKey(privKey)
	if err != nil {
		return fmt.Errorf("Can't Create Signer From Private Key: %w", err)
	}

	// Find IP of bastion
	bastionIP, err := c.FindInstancePublicIP(bastionId)
	if err != nil {
		return fmt.Errorf("Can't Find Public IP on %v: %w", bastionId, err)
	}

	instanceIP, err := c.FindInstancePrivateIP(instanceId)
	if err != nil {
		return fmt.Errorf("Can't Find Public IP on %v: %w", bastionId, err)
	}

	bastionAddress := bastionIP + ":22"
	instanceAddress := instanceIP + ":22"

	sshClientConfig := &ssh.ClientConfig{
		User: username,
		Auth: []ssh.AuthMethod{
			ssh.PublicKeys(signer),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	instanceClient, err := tunnel(instanceAddress, bastionAddress, sshClientConfig)
	if err != nil {
		return err
	}

	return shell(instanceClient)
}

func tunnel(instanceAddress string, bastionAddress string, sshClientConfig *ssh.ClientConfig) (*ssh.Client, error) {
	// https://stackoverflow.com/a/35924799
	bastionClient, err := ssh.Dial("tcp", bastionAddress, sshClientConfig)
	if err != nil {
		return nil, fmt.Errorf("Can't Connect To Bastion: %w", err)
	}

	instanceConn, err := bastionClient.Dial("tcp", instanceAddress)
	if err != nil {
		return nil, fmt.Errorf("Can't Connect To Instance: %w", err)
	}

	ncc, chans, reqs, err := ssh.NewClientConn(instanceConn, instanceAddress, sshClientConfig)
	if err != nil {
		return nil, fmt.Errorf("Can't Connect SSH To Instance: %w", err)
	}

	instanceClient := ssh.NewClient(ncc, chans, reqs)

	return instanceClient, nil
}

func shell(client *ssh.Client) error {
	// Establish an SSH Session
	session, err := client.NewSession()
	if err != nil {
		return fmt.Errorf("Failed to create session: %w\n", err)
	}
	defer session.Close()

	// Set Terminal Into Raw Mode
	fd := int(os.Stdin.Fd())
	oldState, err := terminal.MakeRaw(fd)
	if err != nil {
		panic(err)
	}
	defer terminal.Restore(fd, oldState)

	// Echo required for feedback in raw mode
	modes := ssh.TerminalModes{
		ssh.ECHO:          1,     // enable echoing
		ssh.TTY_OP_ISPEED: 14400, // input speed = 14.4kbaud
		ssh.TTY_OP_OSPEED: 14400, // output speed = 14.4kbaud
	}

	termWidth, termHeight, err := terminal.GetSize(fd)

	// Request a Pseudo Terminal
	if err := session.RequestPty("xterm", termHeight, termWidth, modes); err != nil {
		session.Close()
		return fmt.Errorf("request for pseudo terminal failed: %s", err)
	}

	// Connect All Your IO
	if err := setupIO(session); err != nil {
		return err
	}

	// Start a shell
	if err := session.Shell(); err != nil {
		return fmt.Errorf("failed to start shell: %s", err)
	}

	// Wait for the shell to exit
	session.Wait()

	return nil
}

func setupIO(sess *ssh.Session) error {
	stdin, err := sess.StdinPipe()
	if err != nil {
		return fmt.Errorf("Unable to setup stdin for session: %v", err)
	}
	go io.Copy(stdin, os.Stdin)

	stdout, err := sess.StdoutPipe()
	if err != nil {
		return fmt.Errorf("Unable to setup stdout for session: %v", err)
	}
	go io.Copy(os.Stdout, stdout)

	stderr, err := sess.StderrPipe()
	if err != nil {
		return fmt.Errorf("Unable to setup stderr for session: %v", err)
	}
	go io.Copy(os.Stderr, stderr)
	return nil
}
