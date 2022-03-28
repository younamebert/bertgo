package lexer

// TokenType ast细分后的
type TokenType int

type Token struct {
	Typ   TokenType
	Value string
}
