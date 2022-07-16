package ops

import (
	"github.com/glennliao/task/tasker/op"
	"github.com/glennliao/task/tasker/ops/internal/glob"
)

func init() {
	op.AddOp(op.Op{
		Name: "glob",
		Handler: func(ctx op.Context, args []string) {
			list, _ := glob.Glob(glob.Pattern(args[0]))
			ctx.SetRetVal(ctx.GetVM().ToValue(list))
		},
	})
}
