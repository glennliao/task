package op

import (
	"bufio"
	"fmt"
	"github.com/fatih/color"
	"golang.org/x/text/encoding/simplifiedchinese"
	"io"
	"os/exec"
	"runtime"
	"sync"
)

func init() {
	AddOp(Op{
		Name: "cmd",
		Handler: func(args []string) {
			color.Green("# %s %v\n", "cmd", args[0])
			Cmd(args)
		},
	})
}

// https://developer.aliyun.com/article/934186

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
			fmt.Print(ConvertByte2String([]byte(readString), GB18030))
		}
	}()

	err = cmd.Start()

	err = cmd.Wait()
	if err != nil {
		panic(err)
	}

	wg.Wait()

	if err != nil {
		panic(err)
	}
}

type Charset string

const (
	UTF8    = Charset("UTF-8")
	GB18030 = Charset("GB18030")
)

func ConvertByte2String(byte []byte, charset Charset) string {
	var str string
	switch charset {
	case GB18030:
		var decodeBytes, _ = simplifiedchinese.GB18030.NewDecoder().Bytes(byte)
		str = string(decodeBytes)
	case UTF8:
		fallthrough
	default:
		str = string(byte)
	}
	return str
}
