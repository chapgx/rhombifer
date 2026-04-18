// Package tokens handlers the tokens from the input text parser
package tokens

type Type string

type Token struct {
	Type    Type
	Literal string
}

// NOTE: I wonder if i can encode string against int type in the token type or should I just let the user
// handle this

const (
	ILLEGAL     = "ILLEGAL"
	EOF         = "EOF"
	DASH        = "DASH"
	DOUBLE_DASH = "DOUBLE_DASH"
	IDENT       = "IDENT" // values
	COMMAND     = "COMMAND"
	FLAG        = "FLAG"
	QUOTE       = "\""
)

var (
	flags    map[string]struct{} // user define flags
	commands map[string]struct{} // user define commands
)

// AddFlagKeyWord adds the flagname as a keyword
func AddFlagKeyWord(flagname string) {
	if flags == nil {
		flags = make(map[string]struct{})
	}

	flags[flagname] = struct{}{}
}

// AddCommandKeyWord adds the commandname as a keyword
func AddCommandKeyWord(commandname string) {
	if commands == nil {
		commands = make(map[string]struct{})
	}

	commands[commandname] = struct{}{}
}

// ResloveKeyWord takes an ident and resolves it to a keyword token if any
// otherwise it returns an identifier token
func ResloveKeyWord(ident string) Token {
	_, ok := commands[ident]
	if ok {
		return Token{COMMAND, ident}
	}

	_, ok = flags[ident]
	if ok {
		return Token{FLAG, ident}
	}

	return Token{IDENT, ident}
}

func TokenFromType(t Type) Token {
	return Token{t, string(t)}
}
