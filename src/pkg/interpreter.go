package jlox

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
		}
	}

	return nil, nil
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
