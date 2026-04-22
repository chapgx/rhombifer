package tests

import (
	"fmt"
	"os"
	"testing"

	rhombi "github.com/racg0092/rhombifer"
	"github.com/racg0092/rhombifer/pkg/builtin"
	"github.com/racg0092/rhombifer/pkg/models"
	"github.com/racg0092/rhombifer/tokens"
)

func TestRootAndExe(t *testing.T) {
	root := rhombi.Root()
	root.Run = func(args ...string) error {
		fmt.Println("hello world from root")
		return nil
	}

	t.Run("running root with no values", func(t *testing.T) {
		tokens.RegisterCommand("root")
		os.Args = mimicOsArgs("")

		if err := rhombi.Start(); err != nil {
			t.Error(err)
		}
	})

	t.Run("running root wiht not values and help as default", func(t *testing.T) {
		os.Args = mimicOsArgs("")
		rhombi.GetConfig().RunHelpIfNoInput = true
		help := builtin.HelpCommand(nil, nil)
		root.AddSub(&help)
		if err := rhombi.Start(); err != nil {
			t.Error(err)
		}
	})

	t.Run("running root with flag", func(t *testing.T) {
		os.Args = mimicOsArgs("--lol")

		flag := models.Flag{Name: "lol", Short: "Lol command"}
		rhombi.GetConfig().RunHelpIfNoInput = true
		help := builtin.HelpCommand(nil, nil)
		root.AddSub(&help)
		root.AddFlags(&flag)
		root.Run = func(a ...string) error {
			fmt.Println("Yay from root")
			return nil
		}

		if err := rhombi.Start(); err != nil {
			t.Error(err)
		}
	})

	t.Run("running root with flag and values", func(t *testing.T) {
		os.Args = mimicOsArgs("-r foo bar")

		flag := models.Flag{Name: "run", ShortFormat: "r"}
		rhombi.GetConfig().RunHelpIfNoInput = true
		help := builtin.HelpCommand(nil, nil)
		root.AddSub(&help)
		root.AddFlags(&flag)

		root.Run = func(a ...string) error {
			fmt.Println("Yay from root")
			return nil
		}
		if err := rhombi.Start(); err != nil {
			t.Error(err)
		}
	})
}
