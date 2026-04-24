package rhombifer

import (
	"fmt"
	"os"
	"strings"
	"text/tabwriter"

	text "github.com/chapgx/rhombifer/pkg/text"
)

var helpcmd *Command

// HelpCommand a default built in `help` command.
func HelpCommand() *Command {
	helpcmd = &Command{
		Name:      "help",
		ShortDesc: "Displays help information",
		LongDesc: `
		Displays help information for the specified command or the root command if no command is specified.
		`,
		Leaf: true,
		Run: func(args ...string) error {
			cmd, _ := root.CheckSubCommand("help")
			if len(cmd.Values) == 0 {
				return nil
			}

			sub := root
			for _, cmdname := range cmd.Values {
				s, e := sub.CheckSubCommand(cmdname)
				if e != nil {
					return e
				}
				sub = s
			}

			subHelp(sub)
			cmd.Values = nil
			return nil
		},
	}

	return helpcmd
}

// subHelp is helper to print a [Command] information to the screen
func subHelp(cmd *Command) {
	fmt.Print("\n")
	fmt.Printf("%s\n", strings.ToUpper(string(cmd.Name[0]))+cmd.Name[1:])
	if cmd.LongDesc != "" {
		fmt.Printf("\n%s\n\n", cmd.LongDesc)
	} else {
		fmt.Printf("\n%s\n\n", cmd.ShortDesc)
	}
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintf(w, "%v", text.Bold("Flags:"))

	if cmd.Subs != nil {
		w = tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
		fmt.Fprintf(w, "\n%v", text.Bold("Commands"))
		for _, sub := range cmd.Subs {
			fmt.Fprintf(w, "\n\t%s\t%s", sub.Name, sub.ShortDesc)
		}
		fmt.Fprintf(w, "\n\n")
		w.Flush()
	}

	if cmd.Flags != nil {
		for _, f := range cmd.Flags {
			if f.ShortFormat != "" {
				fmt.Fprintf(w, "\n\t--%s\t-%s\t%s", f.Name, f.ShortFormat, f.Short)
				continue
			}
			fmt.Fprintf(w, "\n\t--%s\t\t%s", f.Name, f.Short)
		}
		fmt.Fprintf(w, "\n")
		w.Flush()
	}
}
