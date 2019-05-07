package endpoint

import (
	"fmt"

	"golang.org/x/crypto/ssh"
)

// Endpoint represent connection interface infos
type Endpoint struct {
	Host   string
	Port   int
	Config *ssh.ClientConfig
}

func (e *Endpoint) String() string {
	return fmt.Sprintf("%s:%d", e.Host, e.Port)
}
