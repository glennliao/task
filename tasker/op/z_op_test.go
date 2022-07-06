package op

import (
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
