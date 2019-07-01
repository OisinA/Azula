package parser

import (
	"azula/lexer"
	"azula/token"
)

const (
	_ int = iota
	lowest
	equals
	lessgreater
	sum
	product
	prefix
	access
	call
	index
)

var precedences = map[token.Type]int{
	token.EQ:       equals,
	token.NOTEQ:    equals,
	token.LT:       lessgreater,
	token.GT:       lessgreater,
	token.PLUS:     sum,
	token.MINUS:    sum,
	token.SLASH:    product,
	token.ASTERISK: product,
	token.LPAREN:   call,
	token.LBRACKET: index,
	token.ACCESS:   access,
}

type Parser struct {
	l *lexer.Lexer
}
