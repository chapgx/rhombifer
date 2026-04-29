package rhombifer

import (
	"fmt"
	"os"
	"strings"

	"github.com/chapgx/rhombifer/ast"
	"github.com/chapgx/rhombifer/lexer"
	"github.com/chapgx/rhombifer/parser"
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
		if helpcmd == nil {
			return fmt.Errorf("help command is nil")
		}
		return helpcmd.Run()
	}

	input := strings.Join(args, " ")
	l := lexer.New(input)
	parse := parser.New(l)
	program := parse.Parse()

	astroot := program.Root

	current_cmd := root
outer:
	for _, node := range astroot {
		switch entity := node.(type) {
		case *ast.Command:
			c, e := current_cmd.CheckSubCommand(entity.Name)
			if e != nil {
				return e
			}

			if c == helpcmd {
				// HACK: we will dig through subcommands here but i think
				// this should be hanlde better in the parsing state
				subcommand := entity.SubCommand
				for subcommand != nil {
					helpcmd.Values = append(helpcmd.Values, subcommand.Name)
					subcommand = subcommand.SubCommand
				}
				// NOTE: could check for flags here and maybe print a more targeted help snippet
				current_cmd = helpcmd
				break outer
			}

			current_cmd = c
			for _, f := range entity.Flags {
				flag := current_cmd.CheckForFlag(f.Name)
				if flag == nil {
					return fmt.Errorf("flag %s was not found in command %s", f.Name, current_cmd.Name)
				}
				foundflags = append(foundflags, flag)

				if !flag.RequiresValue {
					current_cmd.Values = append(current_cmd.Values, f.Value...)
				} else {
					flag.Values = append(flag.Values, f.Value...)
				}
			}

			for _, v := range entity.Values {
				current_cmd.Values = append(current_cmd.Values, v.Content)
			}

			if entity.SubCommand != nil {
				cmd, e := evaluate_command(current_cmd, entity.SubCommand)
				if e != nil {
					return e
				}
				current_cmd = cmd
			}

			continue
		case *ast.Flag:
			flag := current_cmd.CheckForFlag(entity.Name)
			if flag == nil {
				return fmt.Errorf("flag %s was not found in command %s", entity.Name, current_cmd.Name)
			}
			foundflags = append(foundflags, flag)

			if !flag.RequiresValue {
				current_cmd.Values = append(current_cmd.Values, entity.Value...)
			} else {
				flag.Values = append(flag.Values, entity.Value...)
			}
		case *ast.Value:
			current_cmd.Values = append(root.Values, entity.Content)
		}
	}

	return current_cmd.Run(args...)
}

// evaluate_command recursivley evaluates astcmd against cmd
func evaluate_command(cmd *Command, astcmd *ast.Command) (*Command, error) {
	cmd, e := cmd.CheckSubCommand(astcmd.Name)
	if e != nil {
		return nil, fmt.Errorf("command %s was not found", astcmd.Name)
	}

	for _, f := range astcmd.Flags {
		flag := cmd.CheckForFlag(f.Name)
		if flag == nil {
			return nil, fmt.Errorf("flag %s was not found in command %s", f.Name, cmd.Name)
		}
		foundflags = append(foundflags, flag)

		if !flag.RequiresValue {
			cmd.Values = append(cmd.Values, f.Value...)
		} else {
			flag.Values = append(flag.Values, f.Value...)
		}
	}

	for _, v := range astcmd.Values {
		cmd.Values = append(cmd.Values, v.Content)
	}

	if astcmd.SubCommand != nil {
		return evaluate_command(cmd, astcmd.SubCommand)
	}

	return cmd, nil
}
