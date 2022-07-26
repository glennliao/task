package op

import (
	"github.com/fatih/color"
	"log"
	"testing"
)

//func TestTar(t *testing.T) {
//	Tar("dist.tar.gz", "./task")
//}

func TestCmd(t *testing.T) {
	Cmd([]string{"ping www.baidu.com"})
}

func TestTar(t *testing.T) {
	Tar("./a/test.tar.gz", "./a/test.go")
}

func TestSSH(t *testing.T) {
	_, err := createSSHClient("pi@192.168.31.70", "")
	if err != nil {
		log.Fatal(color.RedString(err.Error()))
	}
}
