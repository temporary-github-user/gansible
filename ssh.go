package main

import (
	"io/ioutil"
	"os"

	"golang.org/x/crypto/ssh"
)

func createSigner(keyname string) (signer ssh.Signer, err error) {
	fp, err := os.Open(keyname)
	if err != nil {
		return
	}
	defer fp.Close()

	buf, _ := ioutil.ReadAll(fp)
	signer, _ = ssh.ParsePrivateKey(buf)
	return
}

func makeKeyring() ssh.AuthMethod {
	signers := []ssh.Signer{}
	keys := []string{os.Getenv("HOME") + "/.ssh/id_rsa", os.Getenv("HOME") + "/.ssh/id_dsa"}

	for _, keyname := range keys {
		signer, err := createSigner(keyname)
		if err == nil {
			signers = append(signers, signer)
		}
	}

	return ssh.PublicKeys(signers...)
}
