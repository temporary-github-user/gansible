# Goal

> Run user-selected command on many servers (user-provided as param) with ssh in parallel, collect output from all nodes. The script should print collected output from all nodes on stdout, w/o using temp files.

# Installation & Usage

The package can installed via `go install`

```bash
go install github.com/temporary-github-user/gansible
```

or simply by downloading the binary from the Github Relase page on the right

```
wget https://github.com/temporary-github-user/gansible/releases/download/1.0/gansible_1.0_linux_amd64
```

---

The app uses the first passed argument as a command and all consequent arguments as the hosts where the command will be executed. See the example.

![](images/example.gif)

---

# Thoughts

The package uses ssh keys for authentications and assumes that the default private key is located under the following path `~/.ssh/id_rsa`.

To begin, we dedicate approximately 80% of our CPU resources with an assumption that each CPU thread can handle up to 3 connections(just to begin with sth).

Further, the concurency is handled through the **semaphore pattern** which controls how many goroutines can be running at the same time.

The output from the remote servers captured in a buffered channel, where once any server finishes its job, the output is written to a channel and read as soon as possible.
