package tests

import (
	"strings"

	rhombi "github.com/racg0092/rhombifer"
	"github.com/racg0092/rhombifer/pkg/models"
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
	r := models.Flag{
		Name:        "recursive",
		ShortFormat: "r",
		Required:    true,
	}
	foo := models.Flag{
		Name:        "foo",
		ShortFormat: "f",
	}
	cmd.AddFlags(&r, &foo)
}
