package tasker

import (
	"fmt"
	"github.com/fatih/color"
	"github.com/glennliao/task/tasker/op"
	"github.com/manifoldco/promptui"
	"log"
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
	color.Blue("「 ------- %v ------- 」\n", task.Name)
	for _, step := range task.Steps {
		if step.Op == "need" {
			taskName := step.Args[0]
			task, ok := t.taskMap[taskName]
			if !ok {
				log.Fatal("task " + taskName + " no exists")
			}
			t.runTask(task)
		} else {
			//color.Cyan(" # %d. %-8s %v", i+1, step.Op, step.Args)
			op.OpsRegisterMap[step.Op].Handler(step.Args)
		}
	}
}
