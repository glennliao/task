package ops

import (
	"bufio"
	"fmt"
	"github.com/dop251/goja"
	"github.com/fatih/color"
	"github.com/glennliao/task/tasker/op"
	"github.com/glennliao/task/tasker/ops/internal/scp"
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
	op.AddOp(op.Op{
		Name: "ssh",
		Handler: func(ctx op.Context, args []string) {
			color.Green("# %s \n", "ssh")
			host, password := args[0], args[1]

			sshClient, err := createSSHClient(host, password)
			if err != nil {
				log.Fatal(color.RedString(err.Error()))
			}

			vm := ctx.GetVM()
			cwd := "~"

			sshClientObj := vm.NewObject()

			sshClientObj.Set("upload", func(call goja.FunctionCall) goja.Value {

				from := call.Argument(0).String()
				to := call.Argument(1).String()

				color.Green("# scp %v -> %v\n", from, to)

				if !util.FileExist(from) {

					log.Fatal(color.RedString("the file is not exists, please check it: %s", from))
				}

				n, err := scp.CopyTo(sshClient, from, to, func(cur, total int64) {
					fmt.Printf("\r upload  %d%%", cur*100/total)
				})

				if err == nil || err == io.EOF {
					color.Green("# scp sent %v", util.FormatFileSize(n))
				} else {

					panic(err)
				}

				return goja.Undefined()
			})

			sshClientObj.Set("cmd", func(call goja.FunctionCall) goja.Value {

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

				outText := ""

				go func() {
					defer wg.Done()
					reader := bufio.NewReader(stdout)
					for {
						readString, err := reader.ReadString('\n')
						if err != nil || err == io.EOF {
							return
						}
						outText += readString
						fmt.Print(util.Convert2Utf8Str(readString))
					}
				}()
				color.Green("# %s %s\n", "ssh", call.Argument(0).String())

				err = session.Start("cd " + cwd + " && " + call.Argument(0).String())

				err = session.Wait()
				if err != nil {
					panic(err)
				}

				wg.Wait()
				if err != nil {
					panic(err)
				}

				return vm.ToValue(outText)
			})

			sshClientObj.Set("cd", func(call goja.FunctionCall) goja.Value {
				color.Green("# %s cd %s\n", "ssh", call.Argument(0).String())
				cwd = call.Argument(0).String()
				return goja.Undefined()
			})

			ctx.SetRetVal(sshClientObj)
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
