package main

import (
	_ "embed"
	"fmt"
	"github.com/glennliao/task/tasker"
)

const VERSION = "1.1.0"

//go:embed tasker.js
var taskerJs string

func main() {

	fmt.Println("=== tasker " + VERSION + " ===")

	t := tasker.Tasker{}
	t.Init(taskerJs)
	ok := t.Load()
	if ok {
		t.Run()
	}
}
