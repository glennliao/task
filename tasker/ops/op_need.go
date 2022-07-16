package ops

import (
	"github.com/glennliao/task/tasker/op"
)

func init() {
	op.AddOp(op.Op{
		Name: "need",
		Handler: func(ctx op.Context, args []string) {
			ctx.RunTask(args[0])
		},
	})
}
