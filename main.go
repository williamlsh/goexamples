// Reference: https://blog.ralch.com/tutorial/golang-ssh-tunneling/
package main

import (
	"fmt"
	"io"
	"net"
	"os"

	"golang.org/x/crypto/ssh/agent"

	"golang.org/x/crypto/ssh"
)

// Endpoint represents any server endpoint.
// There are three kinds of endpoints:
// local server, intermediate server, remote/target server
type Endpoint struct {
	Host string
	Port int
}

func (endpoint *Endpoint) String() string {
	return fmt.Sprintf("%s:%d", endpoint.Host, endpoint.Port)
}

// SSHTunnel is a ssh tunnel like this:
// The client is connecting to local endpoint.
// Then the server endpoint mediates between local endpoint and remote endpoint.
type SSHTunnel struct {
	Local  *Endpoint
	Server *Endpoint
	Remote *Endpoint

	Config *ssh.ClientConfig
}

// Start establishes a local server and forwards requests to the intermediate server.
func (t *SSHTunnel) Start() error {
	l, err := net.Listen("tcp", t.Local.String())
	if err != nil {
		return err
	}
	defer l.Close()

	for {
		conn, err := l.Accept()
		if err != nil {
			return err
		}
		go t.forward(conn)
	}
}

// Port forwarding is processed by establishing an SSH connection to the intermediate server.
func (t *SSHTunnel) forward(localConn net.Conn) {
	serverConn, err := ssh.Dial("tcp", t.Server.String(), t.Config)
	if err != nil {
		fmt.Printf("server dial error: %s\n", err)
		return
	}

	remoteConn, err := serverConn.Dial("tcp", t.Remote.String())
	if err != nil {
		fmt.Printf("remote dail error: %s\n", err)
		return
	}

	copyConn := func(writer, reader net.Conn) {
		_, err := io.Copy(writer, reader)
		if err != nil {
			fmt.Printf("io.Copy error: %s\n", err)
		}
	}

	go copyConn(localConn, remoteConn)
	go copyConn(remoteConn, localConn)
}

func SSHAgent() ssh.AuthMethod {
	if sshAgent, err := net.Dial("unix", os.Getenv("SSH_AUTH_SOCK")); err != nil {
		return ssh.PublicKeysCallback(agent.NewClient(sshAgent).Signers)
	}
	return nil
}

func main() {
	localEndpoint := &Endpoint{
		Host: "localhost",
		Port: 9000,
	}
	serverEndpoint := &Endpoint{
		Host: "example.com",
		Port: 22,
	}
	remoteEndpoint := &Endpoint{
		Host: "localhost",
		Port: 8080,
	}

	sshConfig := &ssh.ClientConfig{
		User: "user_name",
		Auth: []ssh.AuthMethod{
			SSHAgent(),
		},
	}

	tunnel := &SSHTunnel{
		Config: sshConfig,
		Local:  localEndpoint,
		Server: serverEndpoint,
		Remote: remoteEndpoint,
	}

	tunnel.Start()
}
