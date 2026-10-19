package jlox

import (
	"fmt"
	"strings"
)

func parenthesizeAST(name string, exprs ...Expr) string {
	var builder strings.Builder

	builder.WriteString("(")
	builder.WriteString(name)

	for _, expr := range exprs {
		builder.WriteString(" ")
		builder.WriteString(PrintAST(expr))
	}

	builder.WriteString(")")

	return builder.String()
}

func PrintAST(expr Expr) string {
	switch e := expr.(type) {
	case Binary:
		return parenthesizeAST(e.operator.Lexeme, e.left, e.right)
	case Grouping:
		return parenthesizeAST("group", e.expression)
	case Comma:
		return parenthesizeAST("comma", e.Expressions...)
	case Conditional:
		return parenthesizeAST("conditional", e.condition, e.trueBranch, e.falseBranch)
	case Literal:
		switch v := e.value.(type) {
		case nil:
			return "nil"
		case string:
			return fmt.Sprintf("%q", v)
		default:
			return fmt.Sprintf("%v", v)
		}
	case Unary:
		return parenthesizeAST(e.operator.Lexeme, e.right)
	default:
		return "Unknown Expression"
	}
}
