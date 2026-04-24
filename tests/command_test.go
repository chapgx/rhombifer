package tests

import (
	"os"
	"testing"

	"github.com/chapgx/rhombifer"
)

func imflagball(args ...string) error {
	return nil
}

func TestHelpCommand(t *testing.T) {
	root := rhombifer.Root()
	help := rhombifer.HelpCommand()
	foo := &rhombifer.Command{Name: "foo", ShortDesc: "i am foo"}
	bar := &rhombifer.Command{Name: "bar", ShortDesc: "i am bar"}

	root.AddSubs(help, foo)
	foo.AddSubs(bar)

	foo.AddFlags(
		rhombifer.NewFlag("ball", "This is a flag").SetRequired().Exec(imflagball),
	)

	args := []string{"programname", "help", "foo"}
	os.Args = args

	e := rhombifer.Start()
	if e != nil {
		t.Fatal(e)
	}

	os.Args = []string{"programname", "help", "foo", "bar"}
	e = rhombifer.Start()
	if e != nil {
		t.Fatal(e)
	}
}
