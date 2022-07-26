package op

import "os"

func init() {
	AddOp(Op{
		Name: "mv",
		Handler: func(args []string) {
			os.Rename(args[0], args[1])
		},
	})
}
