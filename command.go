package rhombifer

import (
	"errors"
	"fmt"
	"strings"

	"github.com/racg0092/rhombifer/pkg/models"
	"github.com/racg0092/rhombifer/pkg/parsing"
	"github.com/racg0092/rhombifer/tokens"
)

type Run func(args ...string) error

type Command struct {
	// command name
	Name string

	// short description showed when the help command is run
	ShortDesc string

	// Long command description
	LongDesc string

	// flags if any
	Flags []*models.Flag

	// Pointers to required flags if any
	requiredFlags []*models.Flag

	// Sub commands for this command
	Subs map[string]*Command

	// Action perform by the command
	Run Run

	// Signifies if this is the root command
	Root bool

	// Signifies if there are no more commands after this one
	Leaf bool

	Values []string
}

// AddFlags adds a flag to the a command
func (cmd *Command) AddFlags(flags ...*models.Flag) {
	for _, f := range flags {
		if f.Required {
			if cmd.requiredFlags == nil {
				cmd.requiredFlags = make([]*models.Flag, 0)
			}
			cmd.requiredFlags = append(cmd.requiredFlags, f)
		}
		if f.RequiresValue && f.Values == nil {
			f.Values = make([]string, 0)
		}
		cmd.Flags = append(cmd.Flags, f)
		tokens.RegisterFlag(f.Name)
	}
}

// AddSub adds a sub command to a command
func (cmd *Command) AddSub(command *Command) {
	if cmd == nil {
		panic("attempting to set sub command to a nil reference")
	}
	if cmd.Subs == nil {
		cmd.Subs = make(map[string]*Command)
	}
	cmd.Subs[command.Name] = command
	tokens.RegisterCommand(command.Name)
}

// AddSubs adds subs commands to the command
func (cmd *Command) AddSubs(subs ...*Command) {
	if len(subs) > 0 {
		if cmd.Subs == nil {
			cmd.Subs = make(map[string]*Command)
		}
		for _, s := range subs {
			cmd.Subs[s.Name] = s
			tokens.RegisterCommand(s.Name)
		}
	}
}

// ValidateRequiredFlags validates if required flags are found in the input string. If any required flag is missing it returns false
// otherwise true. If no flags are required it returns true.
func (cmd *Command) ValidateRequiredFlags(args []string) bool {
	if len(cmd.requiredFlags) <= 0 {
		return true
	}

	if len(args) == 0 {
		return false
	}

	var missing bool = false
	joinArgs := strings.Join(args, " ")
	for _, f := range cmd.requiredFlags {
		if !strings.Contains(joinArgs, "--"+f.Name) && !strings.Contains(joinArgs, "-"+f.ShortFormat) {
			missing = true
			break
		}
	}

	return !missing
}

// RequiredFlags get required flags
func (cmd *Command) RequiredFlags() *[]*models.Flag {
	return &cmd.requiredFlags
}

// CheckSubCommand check if subcommand exists within the command
func (cmd *Command) CheckSubCommand(subcmd string) (*Command, error) {
	if cmd.Subs == nil {
		return nil, errors.New("no sub sommands set for the command [" + cmd.Name + "]")
	}

	for _, scmd := range cmd.Subs {
		if strings.EqualFold(scmd.Name, subcmd) {
			return scmd, nil
		}
	}

	return nil, errors.New("command " + subcmd + " not found")
}

func (cmd *Command) CheckForFlag(name string) *models.Flag {
	for _, flag := range cmd.Flags {
		if flag.Name == name || flag.ShortFormat == name {
			return flag
		}
	}
	return nil
}

var (
	ErrNoASubCommand = errors.New("invalid command format")
	ErrNoSubCommands = errors.New("no subcommands to look through")
)

// DigThroughSubCommand checks user input looking for sub commands until the last one is found
func DigThroughSubCommand(subcommands map[string]*Command, args []string) (*Command, []string, error) {
	// TODO: Throughout test this logic
	if len(subcommands) <= 0 {
		return nil, args, ErrNoSubCommands
	}

	if len(args) <= 0 {
		return nil, args, ErrNoSubCommandPassed
	}

	sub := args[0]

	validsubcommand := parsing.ValidSubCommand(sub)

	if !validsubcommand {
		return nil, args, ErrNoASubCommand
	}

	nargs := args[1:]

	cmd, exists := subcommands[sub]
	if !exists {
		return nil, args, fmt.Errorf("command %s not found", sub)
	}

	if len(nargs) == 0 {
		return cmd, nargs, nil
	}

	validsubcommand = parsing.ValidSubCommand(nargs[0])

	if validsubcommand {
		return DigThroughSubCommand(cmd.Subs, nargs)
	}

	return cmd, nargs, nil
}
