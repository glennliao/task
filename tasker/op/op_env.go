package op

import (
	"encoding/json"
)

var env map[string]any

func init() {
	AddOp(Op{
		Name: "env",
		Handler: func(args []string) {
			json.Unmarshal([]byte(args[0]), &env)
		},
	})
}
