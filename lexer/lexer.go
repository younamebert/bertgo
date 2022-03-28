package lexer

import (
	"bertgo/lexer/util"
	"bytes"
	"io"
)

//结束符号
const EndToken = "$"

// 词语分析器
type Lexer struct {
	*util.Stream
	endToken string
}

func Analyse(source string) []*Token {
	return NewLexer(bytes.NewBufferString(source), EndToken).Analyse()
}

func NewLexer(r io.Reader, et string) *Lexer {
	s := util.NewStream(r, EndToken)
	return &Lexer{Stream: s, endToken: et}
}

func (l *Lexer) Analyse() []*Token {
	tokens := make([]*Token, 0)
	for l.HasNext() {
		c := l.Next()
		if c == EndToken {
			break
		}
		lookahead := l.Peek()

		if c == " " || c == "\n" || c == "\t" {
			continue
		}

		if "/" == c {
			if lookahead == "/" {
				for l.HasNext() {
					if "\n" == l.Next() {
						break
					}
				}
			} else if lookahead == "*" {
				valid := false
				for l.HasNext() {
					p := l.Next()
					if "*" == p && l.Peek() == "/" {
						l.Next()
						valid = true
						break
					}
				}
				if !valid {
					panic("source comment invalid")
				}
			}

			continue
		}

		if c == "{" || c == "}" || c == "(" || c == ")" {
			tokens = append(tokens, NewToken(BRACKET, c))
			continue
		}

		if c == `"` || c == `'` {
			l.PutBack(c)
			tokens = append(tokens, l.MakeString())
			continue
		}

		if IsLetter(c) {
			l.PutBack(c)
			tokens = append(tokens, l.MakeVarOrKeyword())
			continue
		}
		if IsNumber(c) {
			l.PutBack(c)
			tokens = append(tokens, l.MakeNumber())
			continue
		}

		//+ - .
		//+-: 3+5, +5, 3 * -5
		if (c == "+" || c == "-" || c == ".") && IsNumber(lookahead) {
			var lastToken *Token = nil
			if len(tokens) > 0 {
				lastToken = tokens[len(tokens)-1]
			}

			if nil == lastToken || !lastToken.IsValue() || lastToken.IsOperator() {
				l.PutBack(c)
				tokens = append(tokens, l.MakeNumber())
				continue
			}
		}

		if IsOperator(c) {
			l.PutBack(c)
			tokens = append(tokens, l.MakeOp())
			continue
		}

		panic("unexpected character" + c)
	}

	return tokens
}

func (l *Lexer) MakeString() *Token {
	s := ""
	state := 0
	for l.HasNext() {
		c := l.Next()
		switch state {
		case 0:
			if c == `'` {
				state = 1
			} else {
				state = 2
			}
			s += c
		case 1:
			if `'` == c {
				return NewToken(STRING, s+c)
			} else {
				s += c
			}
		case 2:
			if `"` == c {
				return NewToken(STRING, s+c)
			} else {
				s += c
			}
		}
	}

	panic("make string failed")
}

func (l *Lexer) MakeVarOrKeyword() *Token {
	s := ""
	for l.HasNext() {
		lookahead := l.Peek()
		if IsLiteral(lookahead) {
			s += lookahead
		} else {
			break
		}
		l.Next()
	}

	if IsKeyword(s) {
		return NewToken(KEYWORD, s)
	}

	if "true" == s || "false" == s {
		return NewToken(BOOLEAN, s)
	}

	return NewToken(VARIABLE, s)
}
func (l *Lexer) MakeOp() *Token {
	state := 0

	for l.HasNext() {
		lookahead := l.Next()
		switch state {
		case 0:
			switch lookahead {
			case "+":
				state = 1
			case "-":
				state = 2
			case "*":
				state = 3
			case `/`:
				state = 4
			case `>`:
				state = 5
			case `<`:
				state = 6
			case `=`:
				state = 7
			case `!`:
				state = 8
			case `&`:
				state = 9
			case `|`:
				state = 10
			case `^`:
				state = 11
			case `%`:
				state = 12
			case ",":
				return NewToken(OPERATOR, ",")
			case ";":
				return NewToken(OPERATOR, ";")
			}
		case 1:
			switch lookahead {
			case `+`:
				return NewToken(OPERATOR, "++")
			case `=`:
				return NewToken(OPERATOR, "+=")
			default:
				l.PutBack(lookahead)
				return NewToken(OPERATOR, "+")
			}
		case 2:
			switch lookahead {
			case `-`:
				return NewToken(OPERATOR, "--")
			case `=`:
				return NewToken(OPERATOR, "-=")
			default:
				l.PutBack(lookahead)
				return NewToken(OPERATOR, "-")
			}
		case 3:
			switch lookahead {
			case `=`:
				return NewToken(OPERATOR, "*=")
			default:
				l.PutBack(lookahead)
				return NewToken(OPERATOR, "*")
			}
		case 4:
			switch lookahead {
			case `=`:
				return NewToken(OPERATOR, "/=")
			default:
				l.PutBack(lookahead)
				return NewToken(OPERATOR, "/")
			}
		case 5:
			switch lookahead {
			case `=`:
				return NewToken(OPERATOR, ">=")
			case `>`:
				return NewToken(OPERATOR, ">>")
			default:
				l.PutBack(lookahead)
				return NewToken(OPERATOR, ">")
			}
		case 6:
			switch lookahead {
			case `=`:
				return NewToken(OPERATOR, "<=")
			case `<`:
				return NewToken(OPERATOR, "<<")
			default:
				l.PutBack(lookahead)
				return NewToken(OPERATOR, "<")
			}
		case 7:
			switch lookahead {
			case `=`:
				return NewToken(OPERATOR, "==")
			default:
				l.PutBack(lookahead)
				return NewToken(OPERATOR, "=")
			}
		case 8:
			switch lookahead {
			case `=`:
				return NewToken(OPERATOR, "!=")
			default:
				l.PutBack(lookahead)
				return NewToken(OPERATOR, "!")
			}
		case 9:
			switch lookahead {
			case `&`:
				return NewToken(OPERATOR, "&&")
			case `=`:
				return NewToken(OPERATOR, "&=")
			default:
				l.PutBack(lookahead)
				return NewToken(OPERATOR, "&")
			}
		case 10:
			switch lookahead {
			case `|`:
				return NewToken(OPERATOR, "||")
			case `=`:
				return NewToken(OPERATOR, "|=")
			default:
				l.PutBack(lookahead)
				return NewToken(OPERATOR, "|")
			}
		case 11:
			switch lookahead {
			case `^`:
				return NewToken(OPERATOR, "^^")
			case `=`:
				return NewToken(OPERATOR, "^=")
			default:
				l.PutBack(lookahead)
				return NewToken(OPERATOR, "^")
			}
		case 12:
			switch lookahead {
			case `=`:
				return NewToken(OPERATOR, "%=")
			default:
				l.PutBack(lookahead)
				return NewToken(OPERATOR, "%")
			}
		}
	}

	panic("makeOp failed")
}
