package tests

import (
	// "fmt"
	"fmt"
	"testing"

	"github.com/racg0092/rhombifer"
)

func TestFlagValidation(t *testing.T) {
	t.Run("validating short flag required", func(t *testing.T) {
		osargs := mimicOsArgs("-r")
		cmd := rhombifer.Command{
			Name: "foo",
		}
		addSampleFlags(&cmd)
		valid := cmd.ValidateRequiredFlags(osargs)
		if !valid {
			t.Errorf("expected true but got %v", valid)
		}
	})

	t.Run("validating long flag required", func(t *testing.T) {
		osargs := mimicOsArgs("--recursive")
		cmd := rhombifer.Command{
			Name: "foo",
		}
		addSampleFlags(&cmd)
		valid := cmd.ValidateRequiredFlags(osargs)
		if !valid {
			t.Errorf("expected true but got %v", valid)
		}
	})
}

func TestSubCommand(t *testing.T) {
	root := rhombifer.Root()
	cmd := &rhombifer.Command{
		Name:      "cmd",
		ShortDesc: "",
		Run: func(args ...string) error {
			fmt.Println("im a parent command")
			return nil
		},
	}
	scmd := &rhombifer.Command{
		Name:      "sub",
		ShortDesc: "",
		Run: func(args ...string) error {
			fmt.Println("im a child command")
			return nil
		},
	}

	child := &rhombifer.Command{
		Name: "child",
		Run: func(args ...string) error {
			fmt.Println("im a child of a childk")
			return nil
		},
	}

	child1 := &rhombifer.Command{
		Name: "child1",
		Run: func(args ...string) error {
			fmt.Println("im a child of a child of a child")
			return nil
		},
	}

	child.AddSub(child1)
	scmd.AddSub(child)
	cmd.AddSub(scmd)
	root.AddSub(cmd)
	osargs := []string{"sub", "child", "child1"}

	if err := rhombifer.ExecCommand("cmd", osargs...); err != nil {
		t.Error(err)
	}

	t.Run("no sub specified", func(t *testing.T) {
		if err := rhombifer.ExecCommand("cmd"); err != nil {
			t.Error(err)
		}
	})
}
