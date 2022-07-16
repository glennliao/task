package op

import (
	"context"
	"github.com/dop251/goja"
)

type Op struct {
	Name    string
	Handler func(ctx Context, args []string)
}

type Context interface {
	context.Context
	RunTask(name string)
	GetVM() *goja.Runtime
	GetOriCall() *goja.FunctionCall
	SetRetVal(val goja.Value)
}

var OpsRegisterMap = map[string]Op{}

func AddOp(op Op) {
	OpsRegisterMap[op.Name] = op
}
