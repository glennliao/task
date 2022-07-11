package tasker

import (
	"fmt"
	"github.com/dop251/goja"
	"github.com/glennliao/task/tasker/op"
	"github.com/glennliao/task/tasker/util"
	"github.com/manifoldco/promptui"
	"log"
	"os"
	"path/filepath"
)

func (t *Tasker) Load() bool {

	t.taskMap = map[string]*Task{}

	if !util.FileExist(TaskFileName) {
		prompt := promptui.Select{
			Label:        "taskfile.js not found... , so ?",
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
			f, err := os.Create(TaskFileName)
			if err != nil {
				panic(err)
			}
			f.WriteString("task(\"default\")\n\t.cmd(\"ls\")")
			f.Close()
		}

		return false

	}

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
