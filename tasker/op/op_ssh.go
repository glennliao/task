package op

import (
	"bufio"
	"encoding/json"
	"fmt"
	"github.com/fatih/color"
	"github.com/glennliao/task/tasker/op/internal/scp"
	"github.com/glennliao/task/tasker/util"
	"github.com/mitchellh/go-homedir"
	"golang.org/x/crypto/ssh"
	"io"
	"io/ioutil"
	"log"
	"strings"
	"sync"
)

type SSHStep struct {
	Op   string
	Args []string
}

func init() {
	AddOp(Op{
		Name: "ssh",
		Handler: func(args []string) {
			color.Green("# %s \n", "ssh")
			host, password, _sshSteps := args[0], args[1], args[2]

			sshClient, err := createSSHClient(host, password)
			if err != nil {
				log.Fatal(color.RedString(err.Error()))
			}

			var sshSteps []SSHStep
			json.Unmarshal([]byte(_sshSteps), &sshSteps)
			cwd := "~"
			for _, step := range sshSteps {
				switch step.Op {
				case "upload":

					from := step.Args[0]
					to := step.Args[1]

					color.Green("# scp %v -> %v\n", from, to)

					if !util.FileExist(from) {

						log.Fatal(color.RedString("the file is not exists, please check it: %s", from))
					}

					n, err := scp.CopyTo(sshClient, step.Args[0], step.Args[1], func(cur, total int64) {
						fmt.Printf("\r upload  %d%%", cur*100/total)
					})

					if err == nil || err == io.EOF {
						color.Green("# scp sent %v", util.FormatFileSize(n))
					} else {

						panic(err)
					}

				case "cmd":

					var session *ssh.Session
					session, err = sshClient.NewSession()
					if err != nil {
						panic(err)
					}

					stdout, err := session.StdoutPipe()
					if err != nil {
						panic(err)
					}
					var wg sync.WaitGroup
					wg.Add(1)
					go func() {
						defer wg.Done()
						reader := bufio.NewReader(stdout)
						for {
							readString, err := reader.ReadString('\n')
							if err != nil || err == io.EOF {
								return
							}
							fmt.Print(util.Convert2Utf8Str(readString))
						}
					}()
					color.Green("# %s %s\n", "ssh", step.Args[0])

					err = session.Start("cd " + cwd + " && " + step.Args[0])

					err = session.Wait()
					if err != nil {
						panic(err)
					}

					wg.Wait()
					if err != nil {
						panic(err)
					}

				case "cd":
					color.Green("# %s cd %s\n", "ssh", step.Args[0])
					cwd = step.Args[0]
				}
			}

		},
	})
}

func createSSHClient(host, password string) (*ssh.Client, error) {
	hostSps := strings.Split(host, "@")
	addr := hostSps[1]
	user := hostSps[0]
	if !strings.Contains(addr, ":") {
		addr += ":22"
	}

	var authMethods []ssh.AuthMethod

	if signers := loadKey(); len(signers) > 0 {
		authMethods = append(authMethods, ssh.PublicKeys(signers...))
	}

	if password != "" {
		authMethods = append(authMethods, ssh.Password(password))
	}

	sshCfg := &ssh.ClientConfig{
		Config:            ssh.Config{},
		User:              user,
		Auth:              authMethods,
		HostKeyCallback:   ssh.InsecureIgnoreHostKey(),
		BannerCallback:    nil,
		ClientVersion:     "",
		HostKeyAlgorithms: nil,
		Timeout:           0,
	}

	return ssh.Dial("tcp", addr, sshCfg)
}

func loadKey() (list []ssh.Signer) {
	keys := []string{"~/.ssh/id_ed25519", "~/.ssh/id_rsa"}
	for _, key := range keys {
		keyPath, err := homedir.Expand(key)
		if err != nil {
			log.Fatal("find key's home dir failed", err)
		}
		if util.FileExist(keyPath) {
			key, err := ioutil.ReadFile(keyPath)
			if err != nil {
				log.Fatal("ssh key file read failed", err)
			}

			signer, err := ssh.ParsePrivateKey(key)
			if err != nil {
				log.Fatal("ssh key signer failed", err)
			}
			list = append(list, signer)
		}
	}
	return
}
