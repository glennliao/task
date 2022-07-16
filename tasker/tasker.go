package tasker

import (
	"github.com/dop251/goja"
	_ "github.com/glennliao/task/tasker/ops"
)

const (
	ConfigParentDir        = ".config"
	ConfigDir              = "tasker"
	ConfigUseDir           = "use"
	ConfigTaskerJsFilename = "tasker.js"
	DefaultTask            = "default"
)

type Tasker struct {
	configRoot    string
	configUseRoot string
	taskerJs      string
	taskFile      string
	curTask       *Task
	taskMap       map[string]*Task
	runTaskList   []string
	vm            *goja.Runtime
}

type Step struct {
	Op   string
	Args []string
}

type Task struct {
	Name  string
	Steps []Step
}

type Use struct {
	Name string
	Opt  map[string]any
}
