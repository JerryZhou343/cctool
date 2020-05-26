package app

import (
	"fmt"
	"testing"
)

func TestNewGenerateTask(t *testing.T) {
	task := NewGenerateTask("aa")
	fmt.Println(task)
}
