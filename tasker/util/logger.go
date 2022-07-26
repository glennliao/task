package util

import (
	"fmt"
	"github.com/fatih/color"
	"log"
)

type Logger struct{}

func (*Logger) Error(text string) {
	log.Fatal(color.RedString(text))
}

func (l *Logger) ErrorExit(val any) {
	fmt.Printf("Prompt failed %v\n", val)
	log.Fatal(val)
}
