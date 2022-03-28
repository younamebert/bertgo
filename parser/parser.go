package parser

import "bertgo/parser/ast"

// Parser 解析器
type Parser struct {
	stream *ast.PeekTokenStream
}
