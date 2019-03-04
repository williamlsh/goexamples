// Reference: http://blog.ralch.com/tutorial/golang-ssh-connection/
package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"os"
	"strings"

	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/agent"
)

type sshCommand struct {
	Path   string
	Env    []string
	Stdin  io.Reader
	Stdout io.Writer
	Stderr io.Writer
}

// SSHClient is an ssh client.
type SSHClient struct {
	Config *ssh.ClientConfig
	Host   string
	Port   int
}

func (client *SSHClient) RunCommand(cmd *sshCommand) error {
	var (
		session *ssh.Session
		err     error
	)

	if session, err = client.newSeesion(); err != nil {
		return err
	}
	defer session.Close()

	if err := client.prepareCommand(session, cmd); err != nil {
		return err
	}

	err = session.Run(cmd.Path)
	return err
}

func (client *SSHClient) prepareCommand(session *ssh.Session, cmd *sshCommand) error {
	for _, env := range cmd.Env {
		variable := strings.Split(env, "=")
		if len(variable) != 2 {
			continue
		}

		// transfer some environment variable to the remote machine
		if err := session.Setenv(variable[0], variable[1]); err != nil {
			return err
		}
	}

	// To attach our os.Stdin, os.Stdout and os.Stderr to remote command we should open pipes between the local process and remote process.

	if cmd.Stdin != nil {
		stdin, err := session.StdinPipe()
		if err != nil {
			return fmt.Errorf("Unable to setup stdin for session: %v", err)
		}
		go io.Copy(stdin, cmd.Stdin)
	}

	if cmd.Stdout != nil {
		stdout, err := session.StdoutPipe()
		if err != nil {
			return fmt.Errorf("Unable to setup stdout for session: %v", err)
		}
		go io.Copy(cmd.Stdout, stdout)
	}

	if cmd.Stderr != nil {
		stderr, err := session.StderrPipe()
		if err != nil {
			return fmt.Errorf("Unable to setup stderr for session: %v", err)
		}
		go io.Copy(cmd.Stderr, stderr)
	}

	return nil
}

// newSeesion establishes a new ssh connection for client and opens a new session with this connection.
func (client *SSHClient) newSeesion() (*ssh.Session, error) {
	conn, err := ssh.Dial("tcp", fmt.Sprintf("%s:%d", client.Host, client.Port), client.Config)
	if err != nil {
		return nil, fmt.Errorf("failed to dail: %s", err)
	}

	session, err := conn.NewSession()
	if err != nil {
		return nil, fmt.Errorf("failed to create session: %s", err)
	}

	// Before we will be able to run the command on the remote machine, we should create a pseudo terminal on the remote machine. A pseudoterminal (or “pty”) is a pair of virtual character devices that provide a bidirectional communication channel.
	// We should create an xterm terminal that has 80 columns and 40 rows.
	modes := ssh.TerminalModes{
		ssh.ECHO:          0,     // disabling echoing
		ssh.TTY_OP_ISPEED: 14400, // input speed = 14.4kbaud
		ssh.TTY_OP_OSPEED: 14400, // output speed = 14.4kbaud
	}

	if err := session.RequestPty("xterm", 80, 40, modes); err != nil {
		session.Close()
		return nil, fmt.Errorf("request for pseudo terminal failed: %s", err)
	}

	return session, nil
}

// PublicKeyFile parses private key file.
// It's an option used by authentication by using SSH certificate.
func PublicKeyFile(file string) ssh.AuthMethod {
	buf, err := ioutil.ReadFile(file)
	if err != nil {
		return nil
	}

	key, err := ssh.ParsePrivateKey(buf)
	if err != nil {
		return nil
	}

	return ssh.PublicKeys(key)
}

// SSHAgent obtains all stored private keys via SSH_AUTH_SOCK environment variable which stores the SSH agent unix socket.
// It's an option used by authentication by using SSH certificate.
func SSHAgent() ssh.AuthMethod {
	if sshAgent, err := net.Dial("unix", os.Getenv("SSH_AUTH_SOCK")); err == nil {
		return ssh.PublicKeysCallback(agent.NewClient(sshAgent).Signers)
	}
	return nil
}

func main() {
	// Authentication option using a username and password credentials.
	sshConfig := &ssh.ClientConfig{
		User: "user_name",
		Auth: []ssh.AuthMethod{
			// ssh.Password("password"),
			SSHAgent(),
		},
	}

	client := &SSHClient{
		Config: sshConfig,
		Host:   "example.com",
		Port:   22,
	}

	cmd := &sshCommand{
		Path:   "ls -l $LC_DIR",
		Env:    []string{"LC_DIR=/"},
		Stdin:  os.Stdin,
		Stdout: os.Stdout,
		Stderr: os.Stderr,
	}

	fmt.Printf("Running command: %s\n", cmd.Path)
	if err := client.RunCommand(cmd); err != nil {
		fmt.Fprintf(os.Stderr, "command run error: %s\n", err)
		os.Exit(1)
	}
}
