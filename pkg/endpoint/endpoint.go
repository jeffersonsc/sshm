package endpoint

import "fmt"

// Endpoint represent connection interface infos
type Endpoint struct {
	Host string
	Port int
}

func (e *Endpoint) String() string {
	return fmt.Sprintf("%s:%d", e.Host, e.Port)
}
