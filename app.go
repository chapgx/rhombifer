package rhombifer

import (
	"fmt"
	"os"
	"strings"

	"github.com/racg0092/rhombifer/ast"
	"github.com/racg0092/rhombifer/lexer"
	"github.com/racg0092/rhombifer/parser"
)

type App struct{}

// var apponce sync.Once

// Start kick starts the CLI with certain expectations. For more flexebility handle the start of the application your self
// This may change if future versions
func Start() error {
	if root == nil {
		return fmt.Errorf("root command expected got nil")
	}

	args := os.Args[1:]

	if len(args) == 0 && config.RunHelpIfNoInput {
		helpcmd, e := root.CheckSubCommand("help")
		if e != nil {
			return e
		}
		helpcmd.Run()
		return nil
	}

	input := strings.Join(args, " ")
	l := lexer.New(input)
	parse := parser.New(l)
	program := parse.Parse()

	astroot := program.Root

	for _, node := range astroot {
		switch entity := node.(type) {
		case *ast.Command:
			e := command_path(entity, args)
			return e
		case *ast.Flag:
			e := flag_path(root, entity)
			if e != nil {
				return e
			}
		case *ast.Value:
			root.Values = append(root.Values, entity.Content)
		}
	}

	return root.Run(args...)
}

func command_path(cmd *ast.Command, rawinput []string) error {
	root := Root()
	sub, e := root.CheckSubCommand(cmd.Name)
	if e != nil {
		return e
	}

	// note: digs util it find the last subcommad in the chain
	astcommand := cmd.SubCommand.(*ast.Command)
	for astcommand != nil {
		s, e := sub.CheckSubCommand(astcommand.Name)
		if e != nil {
			return e
		}
		sub = s
	}

	// note: checks for flags and add values from the ast to the actual flag
	for _, astf := range cmd.Flags {
		flag := sub.CheckForFlag(astf.Name)
		if flag == nil {
			return fmt.Errorf("expected to find flag %s got nil", astf.Name)
		}
		flag.Values = append(flag.Values, astf.Value...)
		foundflags = append(foundflags, flag)
	}

	return sub.Run(rawinput...)
}

func flag_path(cmd *Command, flag *ast.Flag) error {
	f := cmd.CheckForFlag(flag.Name)
	if f == nil {
		return fmt.Errorf("flag %s was not found", flag.Name)
	}
	foundflags = append(foundflags, f)
	return nil
}
