package op

func init() {
	AddOp(Op{
		Name: "need",
		Handler: func(args []string) {
			// just register op
		},
	})
}
