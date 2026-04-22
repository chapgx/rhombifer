package tests

import (
	"strings"

	rhombi "github.com/racg0092/rhombifer"
)

// Sample user input
func mimicOsArgs(params string) []string {
	var input string
	if params == "" {
		input = "./myprogram"
	} else {
		input = "./myprogram " + params
	}
	args := make([]string, 0)
	args = append(args, strings.Split(input, " ")...)
	return args
}

func addSampleFlags(cmd *rhombi.Command) {
	r := rhombi.Flag{
		Name:        "recursive",
		ShortFormat: "r",
		Required:    true,
	}
	foo := rhombi.Flag{
		Name:        "foo",
		ShortFormat: "f",
	}
	cmd.AddFlags(&r, &foo)
}
