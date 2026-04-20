// Package tokens handlers the tokens from the input text parser
package tokens

import "fmt"

type Type string

type Token struct {
	Type    Type
	Literal string
}

func (t Token) IsNil() bool {
	return t.Literal == "" && t.Type == ""
}

func (t *Token) Equals(x *Token) bool {
	return t.Literal == x.Literal && t.Type == x.Type
}

// NOTE: I wonder if i can encode string against int type in the token type or should I just let the user
// handle this

const (
	ILLEGAL = "ILLEGAL"
	EOF     = "EOF"
	DASH    = "-"
	IDENT   = "IDENT" // values
	QUOTE   = "\""
	VALUE   = "VALUE"
	COMMAND = "COMMAND"
	FLAG    = "FLAG"
)

// TokenFromIdent takes an ident and resolves it to a keyword token if any
// otherwise it returns an identifier token
func TokenFromIdent(ident string) Token {
	return Token{IDENT, ident}
}

func TokenFromType(t Type) Token {
	return Token{t, string(t)}
}

func IsToken(t Token, typ Type) bool {
	return t.Type == typ
}

func IsTokenCommand(literal string) bool {
	if commands == nil {
		panic("commands map is nil")
	}

	_, ok := commands[literal]
	return ok
}

func IsTokenFlag(literal string) bool {
	if flags == nil {
		panic("commands map is nil")
	}

	_, ok := flags[literal]
	return ok
}

var (
	commands map[string]struct{}
	flags    map[string]struct{}
)

func RegisterCommand(commandname string) {
	if commands == nil {
		commands = make(map[string]struct{})
	}

	commands[commandname] = struct{}{}
}

func RegisterFlag(flagname string) {
	if flags == nil {
		flags = make(map[string]struct{})
	}

	flags[flagname] = struct{}{}
}

// ChangeTokenType changes an [IDENT] token to it's proper type
// [COMMAND], [FLAG] or [VALUE]. This can panic if [Token.Type]
// is not [IDENT]
func ChangeTokenType(tok Token) Token {
	if tok.Type != IDENT {
		panic(fmt.Sprintf("expected a IDENT type but got %s", tok.Type))
	}

	_, ok := commands[tok.Literal]
	if ok {
		return Token{COMMAND, tok.Literal}
	}

	_, ok = flags[tok.Literal]
	if ok {
		return Token{FLAG, tok.Literal}
	}

	return tok
}
