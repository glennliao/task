package ops

import (
	"github.com/glennliao/task/tasker/op"
	"os"
)

func init() {
	op.AddOp(op.Op{
		Name: "rm",
		Handler: func(ctx op.Context, args []string) {
			os.RemoveAll(args[0])
		},
	})
}
