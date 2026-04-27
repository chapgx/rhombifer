// Package parser
package parser

import (
	"fmt"

	"github.com/chapgx/rhombifer/ast"
	"github.com/chapgx/rhombifer/lexer"
	"github.com/chapgx/rhombifer/tokens"
)

type Parser struct {
	l *lexer.Lexer

	prevToken tokens.Token
	currToken tokens.Token
	nextToken tokens.Token
	Errors    []error // TODO: instead of crashed log errors
}

func (p *Parser) parse_ident() ast.Node {
	if !tokens.IsTokenCommand(p.currToken.Literal) {
		p.currToken.Type = tokens.VALUE
		return &ast.Value{Token: p.currToken, Content: p.currToken.Literal}
	}

	p.currToken.Type = tokens.COMMAND
	command := &ast.Command{Token: p.currToken, Name: p.currToken.Literal}

	p.nexttoken()
	for p.currToken.Type != tokens.EOF {
		switch p.currToken.Type {
		case tokens.IDENT:
			node := p.parse_ident()
			switch ent := node.(type) {
			case *ast.Command:
				command.SubCommand = ent
			case *ast.Value:
				command.Values = append(command.Values, *ent)
			}
		case tokens.DASH:
			flags := p.parse_dash()
			command.Flags = append(command.Flags, flags...)
			continue
		case tokens.QUOTE:
			value := p.parse_quote()
			command.Values = append(command.Values, value)
		}
		p.nexttoken()
	}
	return command
}

func (p *Parser) parse_command() ast.Node {
	if p.isprevtoken(tokens.QUOTE) && p.isnexttoken(tokens.QUOTE) {
		p.currToken.Type = tokens.VALUE
		return &ast.Value{Token: p.currToken, Content: p.currToken.Literal}
	}

	command := &ast.Command{Token: p.currToken, Name: p.currToken.Literal}

	p.nexttoken()
	current_command := command
	for p.currToken.Type != tokens.EOF {
		switch p.currToken.Type {
		case tokens.IDENT:
			p.currToken.Type = tokens.VALUE
			value := ast.Value{Token: p.currToken, Content: p.currToken.Literal}
			current_command.Values = append(current_command.Values, value)
		case tokens.DASH:
			flags := p.parse_dash()
			current_command.Flags = append(current_command.Flags, flags...)
			if p.currToken.Type == tokens.DASH {
				continue
			}
		case tokens.QUOTE:
			if !p.isnexttoken(tokens.IDENT) && !p.isnexttoken(tokens.COMMAND) {
				panic(fmt.Sprintf("expected an identifer or command to follow but got a %+v\n", p.nextToken))
			}
			p.nexttoken()
			p.currToken.Type = tokens.VALUE
			if !p.isnexttoken(tokens.QUOTE) {
				panic(fmt.Sprintf("expected a closing quote got a %+v\n", p.nextToken))
			}
			value := ast.Value{Token: p.currToken, Content: p.currToken.Literal}
			current_command.Values = append(current_command.Values, value)
			p.nexttoken()
		case tokens.COMMAND:
			node := p.parse_command()
			switch ent := node.(type) {
			case *ast.Value:
				current_command.Values = append(command.Values, *ent)
			case *ast.Command:
				current_command.SubCommand = ent
				current_command = current_command.SubCommand
				// NOTE: this one is really uncessary it should alwasy return a command or a value
			case *ast.Flag:
				current_command.Flags = append(current_command.Flags, *ent)
			}
		}
		p.nexttoken()
	}

	return command
}

func (p *Parser) parse_quote() ast.Value {
	value := ast.Value{
		Token: tokens.Token{Type: tokens.VALUE},
	}

	p.nexttoken()

	// NOTE: here we are looking for the next quote
	for p.currToken.Type != tokens.QUOTE {
		if p.currToken.Type != tokens.EOF {
			if value.Content == "" {
				value.Content += p.currToken.Literal
			} else {
				value.Content += " " + p.currToken.Literal
			}
		} else {
			// TODO: handle this better we may not want to crash here and just log the errors
			panic("missing closing quote for open quote")
		}
		p.nexttoken()
	}

	return value
}

func (p *Parser) parse_flag() ast.Flag {
	flag := ast.Flag{
		Token: tokens.Token{Type: tokens.FLAG},
	}

	if p.currToken.Type != tokens.IDENT {
		panic(fmt.Sprintf("expected an identifier as the flag name got %s", p.currToken.Literal))
	}

	if !tokens.IsTokenFlag(p.currToken.Literal) {
		panic(fmt.Sprintf("expected a flag but %s is not a flag", p.currToken.Literal))
	}

	p.currToken.Type = tokens.FLAG
	flag.Token.Literal = p.currToken.Literal
	flag.Name = p.currToken.Literal
	p.nexttoken()

	for p.currToken.Type != tokens.DASH && p.currToken.Type != tokens.EOF {
		// note: if it is a quote we parse the quote
		if p.currToken.Type == tokens.QUOTE {
			quote := p.parse_quote()
			flag.Value = append(flag.Value, quote.Content)
			p.nexttoken()
			continue
		}

		flag.Value = append(flag.Value, p.currToken.Literal)
		p.nexttoken()
	}

	return flag
}

// nexttoken moves to nextoken
func (p *Parser) nexttoken() {
	p.prevToken = p.currToken
	p.currToken = p.nextToken
	p.nextToken = p.l.NextToken()
}

func (p *Parser) isnexttoken(t tokens.Type) bool {
	return p.nextToken.Type == t
}

func (p *Parser) isprevtoken(t tokens.Type) bool {
	return p.prevToken.Type == t
}

func (p *Parser) parse_short_flags() []ast.Flag {
	if p.currToken.Type != tokens.IDENT {
		panic(fmt.Sprintf("expected ident token got %+v", p.currToken))
	}

	var flags []ast.Flag
	for _, c := range p.currToken.Literal {
		flag := ast.Flag{
			Token: tokens.Token{Type: tokens.FLAG, Literal: string(c)},
			Name:  string(c),
		}
		flags = append(flags, flag)
	}

	return flags
}

func (p *Parser) parse_dash() []ast.Flag {
	if p.isnexttoken(tokens.DASH) {
		p.nexttoken()
		p.nexttoken()
		flag := p.parse_flag()
		return []ast.Flag{flag}
	}

	if !p.isnexttoken(tokens.IDENT) && !p.isnexttoken(tokens.QUOTE) && !p.isnexttoken(tokens.EOF) {
		panic(fmt.Sprintf("expected an identifier, quote, or a EOF token got %+v\n", p.nextToken))
	}

	p.nexttoken()
	flags := p.parse_short_flags()
	p.nexttoken()

	// NOTE: values are added to the last flag and lifted by the command
	// if the flag does not allow values
	lastflag := &flags[len(flags)-1]

	for p.currToken.Type != tokens.DASH && p.currToken.Type != tokens.EOF {
		switch p.currToken.Type {
		case tokens.QUOTE:
			value := p.parse_quote()
			lastflag.Value = append(lastflag.Value, value.Content)
		case tokens.IDENT:
			p.currToken.Type = tokens.VALUE
			lastflag.Value = append(lastflag.Value, p.currToken.Literal)
		}
		p.nexttoken()
	}

	return flags
}

func (p *Parser) Parse() ast.Program {
	prog := ast.Program{Root: make([]ast.Node, 0)}

outer:
	for p.currToken.Type != tokens.EOF {

		// fmt.Println(p.currToken)
		var node ast.Node

		switch p.currToken.Type {
		case tokens.IDENT:
			node = p.parse_ident()
		case tokens.DASH:
			flags := p.parse_dash()
			for _, f := range flags {
				prog.Root = append(prog.Root, &f)
			}
			continue outer
		case tokens.COMMAND:
			node = p.parse_command()
		case tokens.QUOTE:
			val := p.parse_quote()
			node = &val
		}

		prog.Root = append(prog.Root, node)
		p.nexttoken()
	}

	return prog
}

// New creates and returns a new [Parser]
func New(l *lexer.Lexer) *Parser {
	p := &Parser{l: l}

	p.nexttoken()
	p.nexttoken()

	return p
}
