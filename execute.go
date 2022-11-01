package main

import (
	"bytes"
	"log"

	"golang.org/x/crypto/ssh"
)

func executeCmd(cmd, hostname string, config *ssh.ClientConfig) string {
	conn, err := ssh.Dial("tcp", hostname+":22", config)
	if err != nil {
		log.Fatal(err)
	}
	session, _ := conn.NewSession()
	defer session.Close()

	var stdoutBuf bytes.Buffer
	session.Stdout = &stdoutBuf
	session.Run(cmd)

	return hostname + ": " + stdoutBuf.String()
}
