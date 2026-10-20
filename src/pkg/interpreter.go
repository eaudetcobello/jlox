package jlox

import "fmt"

func evaluate(expr Expr) (any, error) {
	switch e := expr.(type) {
	case Literal:
		return e.value, nil
	case Grouping:
		return evaluate(e.expression)
	case Unary:
		right, err := evaluate(e.right)
		if err != nil {
			return nil, err
		}

		switch e.operator.TokenType {
		case MINUS:
			n, ok := right.(float64)
			if !ok {
				return nil, &RuntimeError{Token: e.operator, Msg: "Operand must be a number."}
			}

			return -n, nil
		case BANG:
			return !isTruthy(right), nil
		default:
			return nil, &RuntimeError{Token: e.operator, Msg: "Missing operand."}
		}
	case Binary:
		left, err := evaluate(e.left)
		if err != nil {
			return nil, err
		}
		right, err := evaluate(e.right)
		if err != nil {
			return nil, err
		}

		leftFloat64, ok := left.(float64)
		if !ok {
			return nil, &RuntimeError{Token: e.operator, Msg: "Left of binary expression must be a number."}
		}
		rightFloat64, ok := right.(float64)
		if !ok {
			return nil, &RuntimeError{Token: e.operator, Msg: "Right of binary expression must be a number."}
		}

		switch e.operator.TokenType {
		case MINUS:
			return leftFloat64 - rightFloat64, nil
		case SLASH:
			return leftFloat64 / rightFloat64, nil
		case STAR:
			return leftFloat64 * rightFloat64, nil
		case PLUS:
			return leftFloat64 + rightFloat64, nil
		default:
			return nil, fmt.Errorf("evaluate: unhandled binary operator %s", e.operator.Lexeme)
		}
	default:
		return nil, fmt.Errorf("evaluate: unhandled expression %T", expr)
	}
}

func isTruthy(v any) bool {
	switch v := v.(type) {
	case nil:
		return false
	case bool:
		return v
	default:
		return true
	}
}
