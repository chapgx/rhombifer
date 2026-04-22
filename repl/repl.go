package repl

import (
	"bufio"
	"fmt"
	"io"

	"github.com/chapgx/rhombifer/lexer"
	"github.com/chapgx/rhombifer/tokens"
)

const PROMPT = ">> "

func Start(i io.Reader, o io.Writer) {
	scanner := bufio.NewScanner(i)

	for {
		fmt.Fprintf(o, PROMPT)
		scanned := scanner.Scan()
		if !scanned {
			return
		}

		line := scanner.Text()
		l := lexer.New(line)

		for tok := l.NextToken(); tok.Type != tokens.EOF; tok = l.NextToken() {
			fmt.Fprintf(o, "%+v\n", tok)
		}
	}
}
