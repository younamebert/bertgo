package ast

import (
	"bertgo/lexer"
	"tinyscript/parser/ast"
)

type PeekTokenStream struct {
	tokens []*lexer.Token //
}

// Parse代码分析转ast树
func Parse(source string) ast.ASTNode {
	//lexer.Analyse(source) 分析语法
	return NewParser(lexer.Analyse(source)).parse()
}

func NewParser(tokens []*lexer.Token) *Parser {
	return &Parser{stream: ast.NewPeekTokenStream(tokens)}
}
