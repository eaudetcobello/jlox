package ast

import jlox "github.com/eaudetcobello/jlox/pkg"

type Expr interface{}

// https://eli.thegreenplace.net/2018/the-expression-problem-in-go/
type Eval interface{}

type Binary struct {
	left     Expr
	operator jlox.Token
	right    Expr
}

type Grouping struct {
	expression Expr
}

type Literal struct {
	value interface{}
}

type Unary struct {
	operator jlox.Token
	right    Expr
}
