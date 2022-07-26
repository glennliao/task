package main

import (
	"fmt"
	"github.com/lukasholzer/go-glob"
	"testing"
)

func TestName(t *testing.T) {
	list, _ := glob.Glob(glob.Pattern("**/*.js"))
	for _, s := range list {
		fmt.Println(s)
	}
}
