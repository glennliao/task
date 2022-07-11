package op

import (
	"bufio"
	"fmt"
	"github.com/fatih/color"
	"github.com/glennliao/task/tasker/util"
	"io"
	"os/exec"
	"runtime"
	"sync"
)

func init() {
	AddOp(Op{
		Name: "cmd",
		Handler: func(args []string) {
			color.Green("# %s %v\n", "[cmd]", args[0])
			Cmd(args)
		},
	})
}

func Cmd(args []string) {
	var cmd *exec.Cmd
	if runtime.GOOS == "windows" {
		cmd = exec.Command("cmd", "/C", args[0])
	} else {
		cmd = exec.Command("bash", "-c", args[0])
	}

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		panic(err)
	}

	errout, err := cmd.StderrPipe()
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
		panic(err)
	}

	wg.Wait()

	if err != nil {
		panic(err)
	}
}
