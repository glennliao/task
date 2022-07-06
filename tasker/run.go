package tasker

import (
	"github.com/fatih/color"
	"github.com/glennliao/task/tasker/op"
	"log"
	"os"
	"strings"
)

func (t *Tasker) Run() {
	taskName := DefaultTask
	if len(os.Args) > 1 {
		taskName = os.Args[1]
	}

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
			log.Fatal("task no found: " + taskName)
		case 1:
			task = matchPrefixTaskList[0]
		default:
			var nameList []string
			for _, t := range matchPrefixTaskList {
				nameList = append(nameList, t.Name)
			}
			log.Fatalf("many task with prefix %v : %v", taskName, nameList)
		}

	}

	t.runTask(task)
}

func (t *Tasker) runTask(task *Task) {
	color.Blue("# task %v\n", task.Name)
	for _, step := range task.Steps {
		if step.Op == "need" {
			taskName := step.Args[0]
			task, ok := t.taskMap[taskName]
			if !ok {
				log.Fatal("task " + taskName + " no exists")
			}
			t.runTask(task)
		} else {
			op.OpsRegisterMap[step.Op].Handler(step.Args)
		}
	}
}
