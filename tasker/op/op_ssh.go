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

			var sshClient *ssh.Client
			hostSps := strings.Split(host, "@")
			addr := hostSps[1]
			user := hostSps[0]
			if !strings.Contains(addr, ":") {
				addr += ":22"
			}

			sshCfg := &ssh.ClientConfig{
				Config: ssh.Config{},
				User:   user,
				Auth: []ssh.AuthMethod{
					publicKeyAuthFunc("~/.ssh/id_rsa"),
					ssh.Password(password),
				},
				HostKeyCallback:   ssh.InsecureIgnoreHostKey(),
				BannerCallback:    nil,
				ClientVersion:     "",
				HostKeyAlgorithms: nil,
				Timeout:           0,
			}

			sshClient, err := ssh.Dial("tcp", addr, sshCfg)
			if err != nil {
				panic(err)
			}

			var sshSteps []SSHStep
			json.Unmarshal([]byte(_sshSteps), &sshSteps)
			cwd := "~"
			for _, step := range sshSteps {
				switch step.Op {
				case "upload":
					color.Green("# scp %v -> %v\n", step.Args[0], step.Args[1])

					n, err := scp.CopyTo(sshClient, step.Args[0], step.Args[1], func(cur, total int64) {
						fmt.Printf("\r upload  %d%%", cur*100/total)
					})
					fmt.Println()

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
					go func() {
						//defer wg.Done()
						reader := bufio.NewReader(stdout)
						for {
							readString, err := reader.ReadString('\n')
							if err != nil || err == io.EOF {
								return
							}
							//fmt.Print(ConvertByte2String([]byte(readString), GB18030))
							log.Print("[ssh] ", readString)
						}
					}()
					color.Green("# %s %s\n", "ssh", step.Args[0])

					err = session.Start("cd " + cwd + " && " + step.Args[0])

					err = session.Wait()
					if err != nil {
						panic(err)
					}
					//var wg sync.WaitGroup
					//wg.Add(1)
					//
					//wg.Wait()
					//
					//if err != nil {
					//	panic(err)
					//}
					//return nil

				case "cd":
					color.Green("# %s cd %s\n", "ssh", step.Args[0])
					cwd = step.Args[0]
				}
			}

		},
	})
}

func publicKeyAuthFunc(kPath string) ssh.AuthMethod {
	keyPath, err := homedir.Expand(kPath)
	if err != nil {
		log.Fatal("find key's home dir failed", err)
	}
	key, err := ioutil.ReadFile(keyPath)
	if err != nil {
		log.Fatal("ssh key file read failed", err)
	}
	// Create the Signer for this private key.
	signer, err := ssh.ParsePrivateKey(key)
	if err != nil {
		log.Fatal("ssh key signer failed", err)
	}
	return ssh.PublicKeys(signer)
}
