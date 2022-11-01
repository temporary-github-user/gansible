package main

import (
	"fmt"
	"net"
	"os"
	"runtime"
	"sync"
	"time"

	"golang.org/x/crypto/ssh"
)

func main() {

	var (
		cmd     = os.Args[1]
		hosts   = os.Args[2:]
		workers = runtime.NumCPU() * 8 / 10      // dedicate ~80% of the CPU Threads for the application
		sem     = make(chan struct{}, workers*3) // to begin, let's assume that each worker can handle 3 connections (it's only an assumption)
		res     = make(chan string, workers*3)

		timeout = time.After(60 * time.Second) // setting the default timeout
		wg      = new(sync.WaitGroup)
	)

	config := &ssh.ClientConfig{
		User: os.Getenv("LOGNAME"),
		Auth: []ssh.AuthMethod{makeKeyring()},
		HostKeyCallback: ssh.HostKeyCallback(func(host string, remote net.Addr, pubKey ssh.PublicKey) error {
			//
			// what I'm doing here is not recommended way of veryfing server keys(to be more specic, their fingerprints)
			// for now, it's nothing more than a placeholder
			//
			// essentially I'm doing
			//
			// /etc/ansible/ansible.cfg
			// [defaults]
			// host_key_checking = False
			return nil
		}),
	}

	for _, host := range hosts {
		sem <- struct{}{} // blocking creation of a goroutine untill we have enough "workers"
		wg.Add(1)
		go func(host string) {
			defer wg.Done()
			defer func() { <-sem }()

			res <- executeCmd(cmd, host, config)
		}(host)
	}

	for i := 0; i < len(hosts); i++ {
		select {
		case res := <-res:
			fmt.Println(res)
		case <-timeout:
			fmt.Println("Timed out!")
			return
		}
	}
	wg.Wait()
}
