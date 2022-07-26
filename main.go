package main

import (
	_ "embed"
	"github.com/glennliao/task/tasker"
	"github.com/urfave/cli/v2"
	"log"
	"os"
)

const VERSION = "1.1.1"

//go:embed tasker.js
var taskerJs string

var taskFile = "./taskfile.js"
var taskName = ""

func main() {

	app := &cli.App{
		Name:           "task",
		Version:        VERSION,
		Usage:          "run my task",
		DefaultCommand: "default",
		Before: func(context *cli.Context) error {
			f := context.String("taskfile")
			if f != "" {
				taskFile = f
			}
			return nil
		},
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "taskfile",
				Aliases: []string{"f"},
			},
		},
		CommandNotFound: func(context *cli.Context, s string) {
			taskName = s
		},
		Commands: []*cli.Command{
			{
				Name: "default",
				Action: func(context *cli.Context) error {
					taskName = "default"
					return nil
				},
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}

	if taskName != "" {
		t := tasker.Tasker{}
		t.Init(taskerJs)
		ok := t.Load(taskFile)
		if ok {
			t.Run(taskName)
		}
	}

}
