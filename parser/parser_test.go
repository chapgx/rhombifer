package parser

import (
	"fmt"
	"testing"

	"github.com/chapgx/rhombifer/ast"
	"github.com/chapgx/rhombifer/lexer"
	"github.com/chapgx/rhombifer/tokens"
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

		expected_tree := []ast.Node{
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

		for index, expected := range expected_tree {
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

	t.Run("short flag", func(t *testing.T) {
		input := `-r hello world`
		l := lexer.New(input)
		pars := New(l)

		program := pars.Parse()

		// TODO: implement a better check that this
		for _, node := range program.Root {
			fmt.Printf("%+v\n", node)
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
