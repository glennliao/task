package ops

import (
	"encoding/json"
	"github.com/glennliao/task/tasker/op"
)

var env map[string]any

func init() {
	op.AddOp(op.Op{
		Name: "env",
		Handler: func(ctx op.Context, args []string) {
			json.Unmarshal([]byte(args[0]), &env)
		},
	})
}
