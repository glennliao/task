package ops

import (
	"github.com/fatih/color"
	"github.com/glennliao/task/tasker/op"
	"log"
)

func init() {
	op.AddOp(op.Op{
		Name: "abort",
		Handler: func(ctx op.Context, args []string) {
			log.Fatalln(color.RedString(args[0]))
		},
	})
}
