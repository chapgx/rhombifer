// Package lexer handles lexical analysis
package lexer

import (
	"github.com/chapgx/rhombifer/tokens"
)

type Lexer struct {
	input        string
	shorttokens  string
	shortidx     int
	position     int  // current position in input (points to the current char)
	readPosition int  // current reading position in input (after current char)
	char         byte // current char under examination
}

func New(input string) *Lexer {
	l := &Lexer{input: input}
	l.readChar()
	return l
}

// readChar reads the next character
func (l *Lexer) readChar() {
	if l.readPosition >= len(l.input) {
		l.char = 0
	} else {
		l.char = l.input[l.readPosition]
	}
	l.position = l.readPosition
	l.readPosition += 1
}

func (l *Lexer) NextToken() tokens.Token {
	var tok tokens.Token

	l.skipWhiteSpace()

	switch l.char {
	case '-':
		tok = tokens.TokenFromType(tokens.DASH)
	case '"':
		tok = tokens.TokenFromType(tokens.QUOTE)
	case 0:
		tok = tokens.TokenFromType(tokens.EOF)
		return tok
	default:
		if isLetter(l.char) {
			ident := l.readIdentifier()
			tok = tokens.TokenFromIdent(ident)
			return tok
		} else {
			tok = tokens.TokenFromType(tokens.ILLEGAL)
		}
	}

	l.readChar()
	return tok
}

// readIdentifier reads an identifier
func (l *Lexer) readIdentifier() string {
	position := l.position
	for isLetter(l.char) {
		l.readChar()
	}
	return l.input[position:l.position]
}

func (l *Lexer) skipWhiteSpace() {
	for l.char == ' ' || l.char == '\t' || l.char == '\n' || l.char == '\r' {
		l.readChar()
	}
}

// isLetter checks if ch is a letter
func isLetter(ch byte) bool {
	// NOTE: this may need to grow
	return 'a' <= ch && ch <= 'z' || 'A' <= ch && ch <= 'Z' || ch == '.' || ch == ',' || ch == ';' || ch == '-' || ch >= '0' && ch <= '9'
}

// peekchar looks at the next character in the input
func (l *Lexer) peekchar() byte {
	if l.readPosition >= len(l.input) {
		return 0
	} else {
		return l.input[l.readPosition]
	}
}
