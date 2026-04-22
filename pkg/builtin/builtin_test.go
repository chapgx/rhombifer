package builtin

import (
	"os"
	"strings"
	"testing"

	"github.com/racg0092/rhombifer"
	"github.com/racg0092/rhombifer/pkg/models"
	// "github.com/racg0092/rhombifer/pkg/models"
)

func FuggazziSubs() []rhombifer.Command {
	commands := make([]rhombifer.Command, 0)

	rcmd := rhombifer.Command{
		Name:      "recursive",
		LongDesc:  "A very very long description. Of bla bla bla",
		ShortDesc: "A short description",
	}

	commands = append(commands, rcmd)

	return commands
}

func FugazziFlags(cmd *rhombifer.Command) {
	foo := models.Flag{
		Name:        "Foo",
		Short:       "A short description of foo",
		Long:        "A very long long description of the flag",
		ShortFormat: "f",
	}
	bar := models.Flag{
		Name:  "Bar",
		Short: "A short description of bar",
		Long:  "A very very long long description of the flag",
	}
	cmd.AddFlags(&foo, &bar)
}

func OsArgs(expand string) []string {
	program := "myprg"
	if expand != "" {
		program = program + " " + expand
	}
	program = strings.ReplaceAll(program, "  ", " ")
	return strings.Split(program, " ")
}

func TestHelpWithValue(t *testing.T) {
	os.Args = OsArgs("help rcmd")
	root := rhombifer.Root()
	help := HelpCommand()
	root.AddSubs(&help)
	rcmd := FuggazziSubs()[0]
	FugazziFlags(&rcmd)
	root.Subs["rcmd"] = &rcmd
	rhombifer.GetConfig().RunHelpIfNoInput = true
	if err := rhombifer.Start(); err != nil {
		t.Error(err)
	}
}

func ExampleHelpCommand() {
	// get root command if any. Warning calling this function
	// creates root command
	root := rhombifer.Root()
	// add built in. can be added to any command
	help := HelpCommand()
	root.AddSubs(&help)
	// run help command
	root.Subs["help"].Run()
}
