package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"os"

	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/agent"
)

// Endpoint represent connection interface infos
type Endpoint struct {
	Host string
	Port int
}

func (e *Endpoint) String() string {
	return fmt.Sprintf("%s:%d", e.Host, e.Port)
}

// SSHTunnel data infos for connection
type SSHTunnel struct {
	Local  *Endpoint
	Server *Endpoint
	Remote *Endpoint

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

func main() {
	localEndpoint := &Endpoint{
		Host: "localhost",
		Port: 9000,
	}

	serverEndpoint := &Endpoint{
		Host: "xx.xxx.xx.xxx",
		Port: 22,
	}

	remoteEndpoint := &Endpoint{
		Host: "xx.xx.xx.xx",
		Port: 22,
	}

	sshConfig := &ssh.ClientConfig{
		User: "user",
		// Auth: []ssh.AuthMethod{
		// 	SSHAgent(),
		// },
		Auth: []ssh.AuthMethod{
			// PublicKeyFile("/home/jeffersons/Documentos/AWS-PENS/Telecom2018SP.pem"),
			PublicKeyFile("/home/USER/.ssh/KEY"),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	tunnel := &SSHTunnel{
		Config: sshConfig,
		Local:  localEndpoint,
		Server: serverEndpoint,
		Remote: remoteEndpoint,
	}

	tunnel.Start()
}
