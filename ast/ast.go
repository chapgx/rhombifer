// Package ast
package ast

import "github.com/racg0092/rhombifer/tokens"

type Node interface {
	GetTokenLiteral() string
	IsNil() bool
}

type Command struct {
	Token      tokens.Token
	Name       string
	SubCommand Node
	Flags      []Flag
	Values     []Value
}

func (c *Command) IsNil() bool {
	return c.Token.IsNil() && c.Name == "" && c.SubCommand == nil && c.Flags == nil
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

type Program struct {
	// The root of the cli command expression. This can be a subcommand of root,
	// flags and/or values passed directly at the root command
	Root   []Node
	Errors []error
}
