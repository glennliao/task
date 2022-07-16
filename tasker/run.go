package tasker

import (
	"fmt"
	"github.com/fatih/color"
	"github.com/manifoldco/promptui"
	"strings"
)

func (t *Tasker) Run(taskName string) {

	task, ok := t.taskMap[taskName]

	if !ok {

		var matchPrefixTaskList []*Task

		for _, _task := range t.taskMap {
			if strings.HasPrefix(_task.Name, taskName) && _task.Name != DefaultTask {
				matchPrefixTaskList = append(matchPrefixTaskList, _task)
			}
		}

		switch len(matchPrefixTaskList) {
		case 0:
			logger.Error("task no found: " + taskName)

		case 1:
			task = matchPrefixTaskList[0]
		default:
			var nameList []string
			for _, t := range matchPrefixTaskList {
				nameList = append(nameList, t.Name)
			}

			nameList = append(nameList, "[Cancel]")

			prompt := promptui.Select{
				Label:        fmt.Sprintf("many task with prefix %v , choose one ?", taskName),
				Items:        nameList,
				HideHelp:     true,
				HideSelected: true,
			}

			_, result, err := prompt.Run()

			if err != nil {
				logger.ErrorExit(err)
			}

			if result == "[Cancel]" {
				return
			}

			task = t.taskMap[result]

		}

	}

	t.runTask(task)
}

func (t *Tasker) runTask(task *Task) {
	t.runTaskList = append(t.runTaskList, task.Name)
	color.Blue("「 ------- %v ------- 」\n", task.Name)

	t.vm.Set("currentTask", task.Name)

	_, err := t.vm.RunString(t.taskFile)
	if err != nil {
		panic(err)
	}

	t.runTaskList = t.runTaskList[0 : len(t.runTaskList)-1]

}
