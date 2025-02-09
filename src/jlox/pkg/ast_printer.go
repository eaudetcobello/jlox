package jlox

import (
	"fmt"
	"strings"
)

type AstPrinter struct{}

func (p *AstPrinter) Print(expr Expr) string {
	var text string
	switch expr.(type) {
	case Binary:
		text = "Binary"
	case Grouping:
		text = "Grouping"
	case Literal:
		text = "Literal"
	case Unary:
		text = "Unary"
	default:
		return "Unknown Expression"
	}

	return p.parenthesize(text, expr)
}

func (p *AstPrinter) parenthesize(name string, exprs ...Expr) string {
	var builder strings.Builder

	builder.WriteString("(")
	builder.WriteString(name)

	for _, expr := range exprs {
		builder.WriteString(" ")
		builder.WriteString(p.stringify(expr))
	}

	builder.WriteString(")")

	return builder.String()
}

func (p *AstPrinter) stringify(expr Expr) string {
	switch e := expr.(type) {
	case Binary:
		return p.parenthesize(e.operator.Lexeme, e.left, e.right)
	case Grouping:
		return p.parenthesize("group", e.expression)
	case Literal:
		if e.value == nil {
			return "nil"
		}

		return fmt.Sprintf("%v", e.value)
	case Unary:
		return p.parenthesize(e.operator.Lexeme, e.right)
	default:
		return "Unknown Expression"
	}
}
