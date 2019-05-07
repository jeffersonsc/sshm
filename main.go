package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/go-ini/ini"

	"github.com/jeffersonsc/sshm/pkg/endpoint"
	"github.com/jeffersonsc/sshm/pkg/sshtunnel"
	"golang.org/x/crypto/ssh"
)

var cfg *ini.File
var remote = flag.String("r", "remote", "Remote server in config.ini section")

func main() {
	flag.Parse()

	err := loadConfig()
	checkError("loadconfig", err)

	localPort, err := cfg.Section("").Key("local_port").Int()
	checkError("localport", err)

	localEndpoint := &endpoint.Endpoint{
		Host: "localhost",
		Port: localPort,
	}

	_, err = cfg.GetSection("server")
	checkError("getServerSection", err)

	serverPort, err := cfg.Section("server").Key("port").Int()
	checkError("serverEndpoint", err)

	serverEndpoint := &endpoint.Endpoint{
		Host: cfg.Section("server").Key("host").String(),
		Port: serverPort,
	}

	_, err = cfg.GetSection(*remote)
	checkError("getRemoteSection", err)

	remotePort, err := cfg.Section(*remote).Key("port").Int()
	checkError("remoteEndpoint", err)

	remoteEndpoint := &endpoint.Endpoint{
		Host: cfg.Section(*remote).Key("host").String(),
		Port: remotePort,
	}

	authMethod := cfg.Section("").Key("auth_type").String()
	authKeyFile := cfg.Section("").Key("auth_key_file").String()
	authPass := cfg.Section("").Key("auth_key_file").String()

	sshConfig := &ssh.ClientConfig{
		User: cfg.Section("server").Key("user").String(),
		Auth: []ssh.AuthMethod{
			sshtunnel.SelectAuthMethod(authMethod, authKeyFile, authPass),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	tunnel := &sshtunnel.SSHTunnel{
		Config: sshConfig,
		Local:  localEndpoint,
		Server: serverEndpoint,
		Remote: remoteEndpoint,
	}

	log.Printf(`From remove old connection conflits key exec: ssh-keygen -R [localhost]:%d `, localPort)
	log.Printf(`From access you remote server ssh 'ssh user@localhost -p %d' `, localPort)

	log.Fatal(tunnel.Start())
}

func checkError(local string, err error) {
	if err != nil {
		log.Printf("Failed execute %s, ERROR: %s", local, err.Error())
		os.Exit(1)
	}
}

func loadConfig() (err error) {
	if _, err := os.Stat("./config.ini"); os.IsNotExist(err) {
		fmt.Println("Config file is not exists. Generate new file")
		file, err := os.Create("./config.ini")
		if err != nil {
			return err
		}

		defer file.Close()

		file.WriteString(`
; Auth Methods types
; ssh_agent
; public_key => auth_key_file required
; password => auth_password required

local_port = 22
auth_type = "ssh_agent"
;auth_key_file = "/path/"
;auth_password = "abc123"

; Bastion server
[server]
host = "localhost"
port = 22
user = "user"

; Server from connection nickname
[remote]
host = "localhost"
port = 22
user = "user"
		`)

		// Close file from using from ini
		file.Close()
	}

	cfg, err = ini.Load("./config.ini")
	if err != nil {
		return err
	}

	return nil
}
