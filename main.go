package main

import (
	"fmt"

	"golang.org/x/crypto/ssh"
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

func main() {

}
