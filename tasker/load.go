package tasker

import (
	"context"
	"fmt"
	"github.com/dop251/goja"
	"github.com/glennliao/task/tasker/op"
	"github.com/glennliao/task/tasker/util"
	"github.com/manifoldco/promptui"
	"log"
	"os"
	"path/filepath"
)

func (t *Tasker) Load(taskFile string) bool {

	t.taskMap = map[string]*Task{}

	if !util.FileExist(taskFile) {
		prompt := promptui.Select{
			Label:        fmt.Sprintf("%s not found..., so ?", taskFile),
			Items:        []string{"Create It", "Cancel"},
			HideHelp:     true,
			HideSelected: true,
		}

		_, result, err := prompt.Run()

		if err != nil {
			fmt.Printf("Prompt failed %v\n", err)
			return false
		}

		if result == "Create It" {
			f, err := os.Create(taskFile)
			if err != nil {
				panic(err)
			}
			f.WriteString("task(\"default\")\n\t.cmd(\"ls\")")
			f.Close()
		}
		return false
	}

	taskFile, err := util.LoadJsFile(taskFile)

	t.taskFile = fmt.Sprintf("(()=>{%s\n})()", taskFile)

	if err != nil {
		panic(err)
	}

	vm := goja.New()

	t.vm = vm

	t.initJsVm(vm)
	vm.Set("currentTask", "")

	_, err = vm.RunString(t.taskerJs)

	_, err = vm.RunString(taskFile)
	if err != nil {
		panic(err)
	}

	return true
}

func (t *Tasker) initJsVm(vm *goja.Runtime) {
	tasker := t.taskerJsObject(vm)

	tasker.Set("new", func(call goja.FunctionCall) goja.Value {
		t.curTask = new(Task)
		t.curTask.Name = call.Argument(0).String()
		t.taskMap[t.curTask.Name] = t.curTask
		return goja.Undefined()
	})

	tasker.Set("use", func(call goja.FunctionCall) goja.Value {

		useName := call.Argument(0).String()

		jsFilePath := filepath.Join(t.configUseRoot, useName+".js")
		if !util.FileExist(jsFilePath) {
			log.Fatal("oh! [use] file is miss: " + useName + ".js")
		}

		taskFile, err := util.LoadJsFile(jsFilePath)
		if err != nil {
			panic(err)
		}
		vm.GlobalObject().Set("opt", call.Argument(1).Export().(map[string]any))
		_, err = vm.RunString(fmt.Sprintf("(()=>{%s\n})()", taskFile))
		if err != nil {
			panic(err)
		}
		return goja.Undefined()
	})

	vm.Set("_tasker", tasker)

	console := vm.NewObject()
	console.Set("log", func(call goja.FunctionCall) goja.Value {
		var out []any
		for _, argument := range call.Arguments {
			out = append(out, argument.String())
		}
		log.Println(out...)
		return goja.Undefined()
	})
	vm.Set("console", console)
}

func (t *Tasker) taskerJsObject(vm *goja.Runtime) *goja.Object {
	obj := vm.NewObject()
	for _, _op := range op.OpsRegisterMap {
		op := _op
		obj.Set(op.Name, func(call goja.FunctionCall) goja.Value {
			var args []string
			for _, arg := range call.Arguments {
				args = append(args, arg.String())
			}
			retVal := call.This
			ctx := &TaskContext{tasker: t, VM: vm, retVal: retVal, oriCall: &call}
			op.Handler(ctx, args)
			return ctx.retVal
		})
	}

	return obj
}

type TaskContext struct {
	context.Context
	tasker  *Tasker
	VM      *goja.Runtime
	retVal  goja.Value
	oriCall *goja.FunctionCall
}

func (c *TaskContext) RunTask(name string) {
	c.tasker.Run(name)
}

func (c *TaskContext) GetVM() *goja.Runtime {
	return c.VM
}

func (c *TaskContext) SetRetVal(val goja.Value) {
	c.retVal = val
}

func (c *TaskContext) GetOriCall() *goja.FunctionCall {
	return c.oriCall
}
