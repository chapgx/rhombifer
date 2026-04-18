// Package ast handlers the asbtrat syntax tree
package ast

import (
	"fmt"

	"github.com/racg0092/rhombifer/tokens"
)

type Statement interface {
	// TokenLiteral returns the literal token
	TokenLiteral() string
	// TokenType returns the token type
	TokenType() tokens.Type
}

type Identifier struct {
	Token tokens.Token
	Value string
}

func (i Identifier) TokenLiteral() string   { return i.Token.Literal }
func (i Identifier) TokenType() tokens.Type { return i.Token.Type }

type Command struct {
	Token tokens.Token
	Name  string
	Value []Statement
}

func (c Command) TokenLiteral() string   { return c.Token.Literal }
func (c Command) TokenType() tokens.Type { return c.Token.Type }

type Flag struct {
	Token tokens.Token
	Name  string
	Value []Statement
}

func (f Flag) TokenLiteral() string   { return f.Token.Literal }
func (f Flag) TokenType() tokens.Type { return f.Token.Type }

type ShortFLag struct {
	Token tokens.Token
	Name  string
}

func (sf ShortFLag) TokenLiteral() string   { return sf.Token.Literal }
func (sf ShortFLag) TokenType() tokens.Type { return sf.Token.Type }

type Program struct {
	Tree []Statement
}

// PrintAST prints the AST into the standard output
func PrintAST(ident int, nodes ...Statement) {
	tab := ""
	for i := 0; i < ident; i += 1 {
		tab += "\t"
	}

	for _, node := range nodes {
		switch t := node.(type) {
		case Command:
			fmt.Printf("%scommand => %s", tab, t.Name)
			if len(t.Value) <= 0 {
				break
			}
			ident := ident
			ident += 1
			PrintAST(ident, t.Value...)
		case Flag:
		case Identifier:
		case ShortFLag:
		}
	}
}
