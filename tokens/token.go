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
	ILLEGAL = "ILLEGAL"
	EOF     = "EOF"
	DASH    = "-"
	IDENT   = "IDENT" // values
	QUOTE   = "\""
)

// TokenFromIdent takes an ident and resolves it to a keyword token if any
// otherwise it returns an identifier token
func TokenFromIdent(ident string) Token {
	return Token{IDENT, ident}
}

func TokenFromType(t Type) Token {
	return Token{t, string(t)}
}
