package client

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"strings"

	"golang.org/x/crypto/ssh"
)

type SSHClient struct {
	config *ssh.ClientConfig
	client *ssh.Client
}

// NewSSHClientConfig returns client configuration for SSH client
func NewSSHClientConfig(user, keyFile, passworkPhrase string) (*SSHClient, error) {
	publicKeyMenthod, err := publicKey(keyFile, passworkPhrase)

	if err != nil {
		return nil, err
	}

	sshConfig := &ssh.ClientConfig{
		User: user,
		Auth: []ssh.AuthMethod{
			publicKeyMenthod,
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	return &SSHClient{
		config: sshConfig,
	}, nil
}

// Connect connects to a node
func (c *SSHClient) Connect(host, port string) error {
	client, err := ssh.Dial("tcp", fmt.Sprintf("%s:%s", host, port), c.config)

	if err != nil {
		return err
	}

	c.client = client

	return nil
}

// Close close the client connection
func (c *SSHClient) Close() error {
	if c.client != nil {
		return c.client.Close()
	}

	return nil
}

// ExecuteBatch execute bach commands
func (c *SSHClient) ExecuteBatch(commands []string, prependSudo bool) error {
	fmt.Printf("**** Host: %s ****\n", c.client.RemoteAddr().String())
	for _, cmd := range commands {
		fmt.Printf("** Execute command: %s **\n", cmd)
		stdout, stderr, err := c.ExecuteOne(cmd, prependSudo)

		if err != nil {
			return err
		}

		if len(stdout) > 0 {
			fmt.Println(stdout)
		}

		if len(stderr) > 0 {
			fmt.Printf("Error: %s\n", stderr)
		}
		fmt.Println()
	}

	return nil
}

// ExecuteOne executes one command
func (c *SSHClient) ExecuteOne(cmd string, prependSudo bool) (stdout, stderr string, err error) {
	sess, err := c.client.NewSession()

	if err != nil {
		return "", "", err
	}

	defer sess.Close()

	var stdoutBuf, stderrBuf bytes.Buffer
	sess.Stdout = &stdoutBuf
	sess.Stderr = &stderrBuf

	if prependSudo {
		cmd = "sudo " + cmd
	}

	_ = sess.Run(cmd)

	return strings.TrimSuffix(stdoutBuf.String(), "\n"), strings.TrimSuffix(stderrBuf.String(), "\n"), nil

}

func publicKey(keyPath, passwordPhrase string) (ssh.AuthMethod, error) {
	key, err := ioutil.ReadFile(keyPath)
	if err != nil {
		return nil, err
	}

	var signer ssh.Signer

	if passwordPhrase == "" {
		signer, err = ssh.ParsePrivateKey(key)
	} else {
		signer, err = ssh.ParsePrivateKeyWithPassphrase(key, []byte(passwordPhrase))
	}

	if err != nil {
		return nil, err
	}
	return ssh.PublicKeys(signer), nil
}
