// A flexible and simple unopinionated library for cli tools
package rhombifer

import (
	"fmt"
	"sync"

	"github.com/racg0092/rhombifer/pkg/models"
	"github.com/racg0092/rhombifer/pkg/parsing"
)

var root *Command

var once sync.Once

// Takes in a pointer to `cmd` and sets it as the root command. The **root** command is only set once
// for the application runtime. It means that `root` will be set only the first time this funtion is call.
// Use it if you need to define the `root` command before initialization
func SetRoot(cmd *Command) {
	if cmd != nil {
		cmd.Root = true
		if cmd.Subs == nil {
			cmd.Subs = make(map[string]*Command)
		}
		if cmd.Flags == nil {
			cmd.Flags = make([]*models.Flag, 0)
		}
		// for thread safety
		once.Do(func() {
			root = cmd
		})
	}
}

// Returns root command. If root has not been initialized it creates a new empty [Command]
// and returns the pointer
func Root() *Command {
	if root == nil {
		c := Command{}
		SetRoot(&c)
	}
	return root
}

// Runs the root command if it has been set and no help command is set as default run
func runRoot(args ...string) error {
	root := Root()
	if root.Run == nil && len(args) > 0 {
		return ErroNoRootRunFunc
	} else if root.Run == nil {
		return nil
	}

	if len(args) > 0 {
		foundFlags, err := parsing.FlagsLookup(root.Flags, args...)
		if err != nil && err != parsing.ErrFlagsNilOrEmpty {
			return err
		}
		ff = &foundFlags
	}

	return root.Run(args...)
}

//todo: this function will probably need to be refactor for better usability

// Executes command passed in. It expects [root] to be set
func ExecCommand(cmd string, args ...string) error {
	root := Root()
	if root == nil {
		return fmt.Errorf("Expected root command to be set found %v", root)
	}

	//TODO: this needs to allow root to run without flags and passed values
	if len(args) == 0 && cmd == "" || (cmd == "" && IsFirstArgFlag(args[0])) {
		return runRoot(args...)
	}

	subcommand, found := root.Subs[cmd]
	if !found {
		return fmt.Errorf("Command %s was not found", cmd)
	}

	childcommand, args, err := DigThroughSubCommand(subcommand.Subs, args)

	//HACK: implemented a new error to avoid failure when no sub command is used
	// this should be handle different in the future. Added ignore error to allow parent command
	// to be run with flags
	if err != nil && err != ErrNoSubCommands && err != ErrNoSubCommandPassed && err != ErrNoASubCommand {
		return err
	}

	if childcommand != nil {
		subcommand = childcommand
	}

	if subcommand.Run == nil {
		return fmt.Errorf("Sub command %s, does not have a valid function (Run)", subcommand.Name)
	}

	if len(subcommand.requiredFlags) > 0 {
		if len(args) == 0 {
			return fmt.Errorf("This command (%s) requires flags. Please check the commands docs", subcommand.Name)
		}
		valid := subcommand.ValidateRequiredFlags(args)
		if !valid {
			return fmt.Errorf("Command [%s] requires the expected flags but found [%v]", subcommand.Name, args)
		}
	}

	if len(args) > 0 {
		rawsOnly, err := parsing.IsRawValuesOnly(args...)
		if err != nil {
			return err
		}
		if !rawsOnly {
			//BUG: flags lookup is failing to identify short format flags
			foundFlags, err := parsing.FlagsLookup(subcommand.Flags, args...)
			if err != nil {
				return err
			}
			ff = &foundFlags
		}
	}

	return subcommand.Run(args...)
}
