package tests

import (
	"fmt"
	"os"
	"testing"

	rhombi "github.com/chapgx/rhombifer"
	"github.com/chapgx/rhombifer/tokens"
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
		help := rhombi.HelpCommand()
		root.AddSubs(&help)
		if err := rhombi.Start(); err != nil {
			t.Error(err)
		}
	})

	t.Run("running root with flag", func(t *testing.T) {
		os.Args = mimicOsArgs("--lol")

		flag := rhombi.Flag{Name: "lol", Short: "Lol command"}
		rhombi.GetConfig().RunHelpIfNoInput = true
		help := rhombi.HelpCommand()
		root.AddSubs(&help)
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

		flag := rhombi.Flag{Name: "run", ShortFormat: "r"}
		rhombi.GetConfig().RunHelpIfNoInput = true
		help := rhombi.HelpCommand()
		root.AddSubs(&help)
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
