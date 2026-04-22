package lexer

import (
	"testing"

	"github.com/chapgx/rhombifer/tokens"
)

func TestLexer(t *testing.T) {
	t.Run("full input", func(t *testing.T) {
		input := `run --foo --bar "hello world" -fg`
		l := New(input)

		expectedtokens := []tokens.Token{
			tokens.TokenFromIdent("run"),
			tokens.TokenFromType(tokens.DASH),
			tokens.TokenFromType(tokens.DASH),
			tokens.TokenFromIdent("foo"),
			tokens.TokenFromType(tokens.DASH),
			tokens.TokenFromType(tokens.DASH),
			tokens.TokenFromIdent("bar"),
			tokens.TokenFromType(tokens.QUOTE),
			tokens.TokenFromIdent("hello"),
			tokens.TokenFromIdent("world"),
			tokens.TokenFromType(tokens.QUOTE),
			tokens.TokenFromType(tokens.DASH),
			tokens.TokenFromIdent("fg"),
		}

		for _, expected := range expectedtokens {
			got := l.NextToken()

			if !are_tokens_equal(expected, got) {
				t.Fatalf("\nexpected: %+v\n\ngot: %+v\n\n", expected, got)
			}
		}
	})

	t.Run("short flags", func(t *testing.T) {
		input := `-f value market`
		l := New(input)

		expectedtokens := []tokens.Token{
			tokens.TokenFromType(tokens.DASH),
			tokens.TokenFromIdent("f"),
			tokens.TokenFromIdent("value"),
			tokens.TokenFromIdent("market"),
			tokens.TokenFromType(tokens.EOF),
		}

		for _, expected := range expectedtokens {
			got := l.NextToken()

			if !are_tokens_equal(expected, got) {
				t.Fatalf("\nexpected: %+v\n\ngot: %+v\n\n", expected, got)
			}
		}
	})
}

func are_tokens_equal(a, b tokens.Token) bool {
	return a.Literal == b.Literal && a.Type == b.Type
}
