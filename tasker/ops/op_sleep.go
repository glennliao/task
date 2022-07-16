package ops

import (
	"github.com/glennliao/task/tasker/op"
	"strconv"
	"time"
)

func init() {
	op.AddOp(op.Op{
		Name: "sleep",
		Handler: func(ctx op.Context, args []string) {
			num, err := strconv.Atoi(args[0])
			if err != nil {
				panic(err)
			}
			time.Sleep(time.Duration(num) * time.Second)
		},
	})
}
