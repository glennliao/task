package op

import (
	"bufio"
	"encoding/json"
	"fmt"
	"github.com/fatih/color"
	"github.com/glennliao/task/tasker/util"
	"io"
	"os"
	"log"
	"os/exec"
	"runtime"
	"sync"
)

func init() {
	AddOp(Op{
		Name: "cmd",
		Handler: func(args []string) {
			color.Green(" # %s %v\n", "[cmd]", args[0])
			Cmd(args)
		},
	})
}

type CmdOption struct {
	Env map[string]any
}

func Cmd(args []string) {
	var cmd *exec.Cmd
	if runtime.GOOS == "windows" {
		cmd = exec.Command("cmd", "/C", args[0])
	} else {
		cmd = exec.Command("bash", "-c", args[0])
	}

	var cmdOption CmdOption

	if len(args) > 1 {
		json.Unmarshal([]byte(args[1]), &cmdOption)

	}

	var envs = os.Environ()
	for k, v := range env {
		envs = append(envs, fmt.Sprintf("%s=%v", k, v))
	}
	for k, v := range cmdOption.Env {
		envs = append(envs, fmt.Sprintf("%s=%v", k, v))
	}

	cmd.Env = envs

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		panic(err)
	}

	errout, err := cmd.StderrPipe()
	if err != nil {
		panic(err)
	}

	var wg sync.WaitGroup
	wg.Add(2)

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

	go func() {
		defer wg.Done()
		reader := bufio.NewReader(errout)
		for {
			readString, err := reader.ReadString('\n')
			if err != nil || err == io.EOF {
				return
			}
			fmt.Print(util.Convert2Utf8Str(readString))
		}
	}()

	err = cmd.Start()
	if err != nil {
		panic(err)
	}
	err = cmd.Wait()

	if err != nil {

		switch err.(type) {
		case *exec.ExitError:
			e := err.(*exec.ExitError)
			log.Println(e.String())
		default:
			log.Fatal(err)
		}

	}

	wg.Wait()
}
