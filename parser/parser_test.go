package parser

import (
	"fmt"
	"testing"

	"github.com/racg0092/rhombifer/ast"
	"github.com/racg0092/rhombifer/lexer"
	"github.com/racg0092/rhombifer/tokens"
)

func TestParser(t *testing.T) {
	t.Run("no command", func(t *testing.T) {
		tokens.RegisterCommand("run")
		tokens.RegisterCommand("all")
		tokens.RegisterFlag("foo")
		tokens.RegisterFlag("bar")

		input := "--foo hello world --bar"
		l := lexer.New(input)
		pars := New(l)

		expected_program := ast.Program{}
		expected_program.Root = []ast.Node{
			&ast.Flag{
				Token: tokens.Token{Type: tokens.FLAG, Literal: "foo"},
				Name:  "foo",
				Value: []string{"hello", "world"},
			},
			&ast.Flag{
				Token: tokens.Token{Type: tokens.FLAG, Literal: "bar"},
				Name:  "bar",
			},
		}

		program := pars.Parse()

		for index, expected := range expected_program.Root {
			got := program.Root[index]
			if !are_flags_equals(expected, got) {
				t.Fatalf("\nexpected: %+v\n\ngot: %+v\n\n", expected, got)
			}
		}
	})

	t.Run("with command", func(t *testing.T) {
		input := `run all --foo hello world --bar "Todo App"`
		l := lexer.New(input)
		pars := New(l)

		program := pars.Parse()

		// TODO: implement a better check that this
		for _, node := range program.Root {
			fmt.Printf("%+v\n", node)
			cmd, ok := node.(*ast.Command)
			if !ok {
				t.Fatalf("expected a command above")
			}

			if cmd.SubCommand != nil {
				fmt.Printf("%+v\n", cmd.SubCommand)
			}
		}
	})
}

func print() {
}

func are_flags_equals(a, b ast.Node) bool {
	af, ok := a.(*ast.Flag)
	if !ok {
		return false
	}

	bf, ok := a.(*ast.Flag)
	if !ok {
		return false
	}

	return af.Equals(bf)
}

// func are_commands_equals(a, b ast.Node) bool {
// 	ac, ok := a.(*ast.Command)
// 	if !ok {
// 		return false
// 	}
//
// 	bc, ok := b.(*ast.Command)
// 	if !ok {
// 		return false
// 	}
//
// 	if ac.SubCommand.IsNil() && !bc.SubCommand.IsNil() {
// 		return false
// 	} else if !ac.SubCommand.IsNil() && bc.SubCommand.IsNil() {
// 		return false
// 	} else {
//
// 	}
// }
