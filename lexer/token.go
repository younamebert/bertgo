package lexer

// TokenType ast็ปๅๅ็
type TokenType int

type Token struct {
	Typ   TokenType
	Value string
}
