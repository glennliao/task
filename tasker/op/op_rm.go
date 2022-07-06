package op

import "os"

func init() {
	AddOp(Op{
		Name: "rm",
		Handler: func(args []string) {
			os.RemoveAll(args[0])
		},
	})
}
