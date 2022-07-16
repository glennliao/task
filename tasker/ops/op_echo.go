package ops

import (
	"github.com/fatih/color"
	"github.com/glennliao/task/tasker/op"
)

func init() {
	op.AddOp(op.Op{
		Name: "echo",
		Handler: func(ctx op.Context, args []string) {
			color.White("%v\n", args[0])
		},
	})
}
