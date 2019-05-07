package sshtunnel

import (
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"os"

	"github.com/jeffersonsc/sshm/pkg/endpoint"
	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/agent"
)

// SSHTunnel data infos for connection
type SSHTunnel struct {
	Local  *endpoint.Endpoint
	Server *endpoint.Endpoint
	Remote *endpoint.Endpoint

	Config *ssh.ClientConfig
}

func (st *SSHTunnel) forward(localConn net.Conn) {
	serverConn, err := ssh.Dial("tcp", st.Server.String(), st.Config)
	if err != nil {
		fmt.Println("Server dial error: ", err.Error())
		return
	}

	remoteConn, err := serverConn.Dial("tcp", st.Remote.String())
	if err != nil {
		fmt.Println("Remote dial error: ", err.Error())
		return
	}

	copyConn := func(writer, reader net.Conn) {
		defer writer.Close()
		defer reader.Close()
		_, err := io.Copy(writer, reader)
		if err != nil {
			fmt.Println("io.Copy error: ", err.Error())
		}
	}

	go copyConn(localConn, remoteConn)
	go copyConn(remoteConn, localConn)
}

// Start local server client
func (st *SSHTunnel) Start() error {
	listener, err := net.Listen("tcp", st.Local.String())
	if err != nil {
		return err
	}
	defer listener.Close()

	for {
		conn, err := listener.Accept()
		if err != nil {
			return err
		}

		go st.forward(conn)
	}
}

// SelectAuthMethod from auth connection
func SelectAuthMethod(method, keypath, pass string) ssh.AuthMethod {
	switch method {
	case "auth_key_file":
		return PublicKeyFile(keypath)
	case "password":
		return Password(pass)
	default:
		return SSHAgent()
	}
}

// SSHAgent return autehnticable method
func SSHAgent() ssh.AuthMethod {
	if sshAgent, err := net.Dial("unix", os.Getenv("SSH_AUTH_SOCK")); err == nil {
		return ssh.PublicKeysCallback(agent.NewClient(sshAgent).Signers)
	}
	return nil
}

// PublicKeyFile method authentication ssh
func PublicKeyFile(file string) ssh.AuthMethod {
	buffer, err := ioutil.ReadFile(file)
	if err != nil {
		return nil
	}

	key, err := ssh.ParsePrivateKey(buffer)
	if err != nil {
		return nil
	}
	return ssh.PublicKeys(key)
}

// Password method auth
func Password(pass string) ssh.AuthMethod {
	return ssh.Password(pass)
}
