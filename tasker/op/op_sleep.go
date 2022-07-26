package op

import (
	"strconv"
	"time"
)

func init() {
	AddOp(Op{
		Name: "sleep",
		Handler: func(args []string) {
			num, err := strconv.Atoi(args[0])
			if err != nil {
				panic(err)
			}
			time.Sleep(time.Duration(num) * time.Second)
		},
	})
}
