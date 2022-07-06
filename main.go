package main

import (
	_ "embed"
	"fmt"
	"github.com/glennliao/task/tasker"
)

const VERSION = "1.0.0"

//go:embed tasker.js
var taskerJs string

func main() {

	fmt.Println("=== tasker " + VERSION + " ===")

	t := tasker.Tasker{}
	t.Init(taskerJs)
	t.Load()
	t.Run()

}
