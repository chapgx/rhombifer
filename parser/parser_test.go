package parser

import (
	"fmt"
	"testing"

	"github.com/chapgx/rhombifer/ast"
	"github.com/chapgx/rhombifer/lexer"
	"github.com/chapgx/rhombifer/tokens"
)

func TestParser(t *testing.T) {
	tokens.RegisterCommand("run")
	tokens.RegisterCommand("all")
	tokens.RegisterFlag("foo")
	tokens.RegisterFlag("bar")

	t.Run("no command", func(t *testing.T) {
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

		expectedroot := []ast.Node{
			&ast.Command{
				Token: tokens.Token{Type: tokens.COMMAND, Literal: "run"},
				Name:  "run",
				SubCommand: &ast.Command{
					Token: tokens.Token{Type: tokens.COMMAND, Literal: "all"},
					Name:  "all",
					Flags: []ast.Flag{
						{
							Token: tokens.Token{Type: tokens.FLAG, Literal: "foo"},
							Name:  "foo",
							Value: []string{"hello", "world"},
						},
						{
							Token: tokens.Token{Type: tokens.FLAG, Literal: "bar"},
							Name:  "bar",
							Value: []string{"Todo App"},
						},
					},
				},
			},
		}

		if len(expectedroot) != len(program.Root) {
			t.Fatal("expectedroot len does not match progam.Root len")
		}

		var e error
		for index, expected := range expectedroot {
			got := program.Root[index]
			e = are_nodes_equal(expected, got)
			if e != nil {
				t.Fatalf("%s\n\n%s", e, fmt.Sprintf("expected:%+v\n\ngot:%+v\n\n", expected, got))
			}
		}
	})

	t.Run("short flag", func(t *testing.T) {
		input := `-r hello world`
		l := lexer.New(input)
		pars := New(l)

		program := pars.Parse()
		expected := ast.Flag{
			Token: tokens.Token{Type: tokens.FLAG, Literal: "r"},
			Name:  "r",
			Value: []string{"hello", "world"},
		}

		node := program.Root[0]

		got, ok := node.(*ast.Flag)
		if !ok {
			t.Fatalf("expected a flag got %+v\n", got)
		}

		are_flags_equals(got, &expected)
	})

	t.Run("command short flag", func(t *testing.T) {
		input := `run all -f hello world`
		l := lexer.New(input)
		par := New(l)
		program := par.Parse()

		expectedroot := []ast.Node{
			&ast.Command{
				Token: tokens.Token{Type: tokens.COMMAND, Literal: "run"},
				Name:  "run",
				SubCommand: &ast.Command{
					Token: tokens.Token{Type: tokens.COMMAND, Literal: "all"},
					Name:  "all",
					Flags: []ast.Flag{
						{
							Token: tokens.Token{Type: tokens.FLAG, Literal: "f"},
							Name:  "f",
							Value: []string{"hello", "world"},
						},
					},
				},
			},
		}

		if len(expectedroot) != len(program.Root) {
			t.Fatal("expectedroot len and program.Root len does not match")
		}

		var e error
		for index, expected := range expectedroot {
			got := program.Root[index]
			e = are_nodes_equal(expected, got)
			if e != nil {
				t.Errorf("%s\n\n%s", e, fmt.Sprintf("expected:%+v\n\ngot:%+v\n\n", expected, got))
				expectedcmd, ok := expected.(*ast.Command)
				if ok {
					if expectedcmd.SubCommand != nil {
						fmt.Printf("go subcommand:%+v\n\n", expectedcmd.SubCommand)
					}
				}
				command, ok := got.(*ast.Command)
				if ok {
					if command.SubCommand != nil {
						fmt.Printf("go subcommand:%+v\n\n", command.SubCommand)
					}
				}
			}
		}
	})
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

func are_nodes_equal(a, b ast.Node) error {
	if a == nil && b == nil {
		return nil
	}

	if a == nil && b != nil {
		return fmt.Errorf("node a is nil (%v) but node b (%v) is not", a, b)
	}

	if b == nil && a != nil {
		return fmt.Errorf("node a (%v) is not nil but node b (%v) is not", a, b)
	}

	switch aval := a.(type) {
	case *ast.Command:
		bc, ok := b.(*ast.Command)
		if !ok {
			return fmt.Errorf("expected node b to be of type Command")
		}
		if !aval.Equals(bc) {
			return fmt.Errorf("node a and be are not equal. command comparison case")
		}
	case *ast.Flag:
		bflag, ok := b.(*ast.Flag)
		if !ok {
			return fmt.Errorf("expected node b to be of type Flag")
		}
		if !aval.Equals(bflag) {
			return fmt.Errorf("node a and be are not equal. flag comparison case")
		}
	case *ast.Value:
		bv, ok := b.(*ast.Value)
		if !ok {
			return fmt.Errorf("expected node b to be of type Value")
		}
		if !aval.Equals(bv) {
			return fmt.Errorf("node a and be are not equal. value comparison case")
		}
	}

	return nil
}
