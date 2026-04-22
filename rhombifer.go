// Package rhombifer a flexible and simple unopinionated library for cli tools
package rhombifer

import (
	"sync"

	"github.com/racg0092/rhombifer/pkg/models"
)

var root *Command

var once sync.Once

// SetRoot Takes in a pointer to `cmd` and sets it as the root command. The **root** command is only set once
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

// Root returns root command. If root has not been initialized it creates a new empty [Command]
// and returns the pointer
func Root() *Command {
	if root == nil {
		c := Command{}
		SetRoot(&c)
	}
	return root
}
