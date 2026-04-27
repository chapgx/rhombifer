// Package ast
package ast

import "github.com/chapgx/rhombifer/tokens"

type Node interface {
	GetTokenLiteral() string
	IsNil() bool
}

type Command struct {
	Token      tokens.Token
	Name       string
	SubCommand *Command
	Flags      []Flag
	Values     []Value
}

func (c *Command) IsNil() bool {
	return c.Token.IsNil() && c.Name == "" && c.SubCommand == nil && c.Flags == nil
}

// Equals checks if cmd instance is equal to x
func (c *Command) Equals(x *Command) bool {
	// TODO: finish implementation i already staed in the test file just copy it over
	// and finish it
	if x == nil {
		return false
	}

	if c.Name != x.Name || !c.Token.Equals(&x.Token) {
		return false
	}

	if len(c.Flags) != len(x.Flags) {
		return false
	}

	for index, f := range c.Flags {
		xf := x.Flags[index]
		if !f.Equals(&xf) {
			return false
		}
	}

	if len(c.Values) != len(x.Values) {
		return false
	}

	for index, v := range c.Values {
		xv := x.Values[index]
		if !v.Equals(&xv) {
			return false
		}
	}

	if c.SubCommand != nil {
		if !c.SubCommand.Equals(x.SubCommand) {
			return false
		}
	} else if x.SubCommand != nil {
		return false
	}

	return true
}

func (c *Command) GetTokenLiteral() string { return c.Token.Literal }

type Flag struct {
	Token tokens.Token
	Name  string
	Value []string
}

func (f *Flag) IsNil() bool {
	return f.Token.IsNil() && f.Name == "" && f.Value == nil
}

func (f *Flag) Equals(x *Flag) bool {
	if !f.Token.Equals(&x.Token) || f.Name != x.Name {
		return false
	}

	for index, v := range f.Value {
		if v != x.Value[index] {
			return false
		}
	}

	return true
}

func (f *Flag) GetTokenLiteral() string { return f.Token.Literal }

type Value struct {
	Token   tokens.Token
	Content string
}

func (v *Value) IsNil() bool             { return v.Content == "" }
func (v *Value) GetTokenLiteral() string { return v.Token.Literal }

// Equals compares v to x and returns true if they are equal
func (v *Value) Equals(x *Value) bool {
	are_equals := v.Token.Equals(&x.Token)
	if !are_equals {
		return false
	}

	return v.Content != x.Content
}

type Program struct {
	// The root of the cli command expression. This can be a subcommand of root,
	// flags and/or values passed directly at the root command
	Root   []Node
	Errors []error
}
