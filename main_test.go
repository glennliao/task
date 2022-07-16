package main

import (
	"github.com/glennliao/task/tasker"
	"testing"
)

func TestHi(tx *testing.T) {
	t := tasker.Tasker{}
	t.Init(taskerJs, taskerDTS)
	ok := t.Load(taskFile)
	taskName = "tar"
	if ok {
		t.Run(taskName)
	}
}
