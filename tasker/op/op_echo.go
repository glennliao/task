package op

import (
	"github.com/fatih/color"
)

func init() {
	AddOp(Op{
		Name: "echo",
		Handler: func(args []string) {
			color.Green("# %v\n", args[0])
		},
	})
}
