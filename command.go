package rhombifer

import (
	"errors"
	"strings"

	"github.com/chapgx/rhombifer/tokens"
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
	Flags []*Flag

	// Pointers to required flags if any
	requiredFlags []*Flag

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
func (cmd *Command) AddFlags(flags ...*Flag) {
	for _, f := range flags {
		if f.Required {
			if cmd.requiredFlags == nil {
				cmd.requiredFlags = make([]*Flag, 0)
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
func (cmd *Command) RequiredFlags() *[]*Flag {
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

func (cmd *Command) CheckForFlag(name string) *Flag {
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
