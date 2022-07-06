package op

var OpsRegisterMap = map[string]Op{}

type Op struct {
	Name    string
	Handler func(args []string)
}

func AddOp(op Op) {
	OpsRegisterMap[op.Name] = op
}
