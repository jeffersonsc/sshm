package main

import (
	"github.com/jeffersonsc/sshm/pkg/endpoint"
	"github.com/jeffersonsc/sshm/pkg/sshtunnel"
	"golang.org/x/crypto/ssh"
)

func main() {
	localEndpoint := &endpoint.Endpoint{
		Host: "localhost",
		Port: 9000,
	}

	serverEndpoint := &endpoint.Endpoint{
		Host: "localhost",
		Port: 22,
	}

	remoteEndpoint := &endpoint.Endpoint{
		Host: "localhost",
		Port: 22,
	}

	sshConfig := &ssh.ClientConfig{
		User: "user",
		Auth: []ssh.AuthMethod{
			sshtunnel.SSHAgent(),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	tunnel := &sshtunnel.SSHTunnel{
		Config: sshConfig,
		Local:  localEndpoint,
		Server: serverEndpoint,
		Remote: remoteEndpoint,
	}

	tunnel.Start()
}
