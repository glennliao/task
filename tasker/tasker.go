package tasker

const (
	ConfigParentDir        = ".config"
	ConfigDir              = "tasker"
	ConfigUseDir           = "use"
	ConfigTaskerJsFilename = "tasker.js"
	TaskFileName           = "./taskfile.js"
	DefaultTask            = "default"
)

type Tasker struct {
	configRoot    string
	configUseRoot string
	taskerJs      string
	curTask       *Task
	taskMap       map[string]*Task
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
