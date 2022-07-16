package ops

import (
	"github.com/glennliao/task/tasker/op"
	"os"
)

func init() {
	op.AddOp(op.Op{
		Name: "mv",
		Handler: func(ctx op.Context, args []string) {
			os.Rename(args[0], args[1])
		},
	})
}
