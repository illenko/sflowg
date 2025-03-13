package sflowg

import "github.com/expr-lang/expr"

func Eval(expression string, context map[string]any) (any, error) {
	return expr.Eval(FormatExpression(expression), context)
}
