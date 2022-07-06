package tasker

import (
	"github.com/dop251/goja"
	"github.com/glennliao/task/tasker/op"
	"github.com/glennliao/task/tasker/util"
	"path/filepath"
)

func (t *Tasker) Load() {

	t.taskMap = map[string]*Task{}

	taskFile, err := util.LoadJsFile(TaskFileName)

	if err != nil {
		panic(err)
	}

	vm := goja.New()
	t.initJsVm(vm)

	_, err = vm.RunString(t.taskerJs)

	_, err = vm.RunString(taskFile)
	if err != nil {
		panic(err)
	}
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

		taskFile, err := util.LoadJsFile(filepath.Join(t.configUseRoot, useName+".js"))
		if err != nil {
			panic(err)
		}
		vm.GlobalObject().Set("opt", call.Argument(1).Export().(map[string]any))
		_, err = vm.RunString(taskFile)
		if err != nil {
			panic(err)
		}

		return goja.Undefined()
	})

	vm.Set("_tasker", tasker)
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
			t.curTask.Steps = append(t.curTask.Steps, Step{
				Op:   op.Name,
				Args: args,
			})

			return call.This
		})
	}

	return obj
}
