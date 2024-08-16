package lox

import (
	"strings"
)

func AstPrinter(name string, exprs ...Expr) string {
	var builder strings.Builder

	builder.WriteString("(")
	builder.WriteString(name)

	for _, expr := range exprs {
		builder.WriteString(" ")
		builder.WriteString(expr.AstPrint())
	}

	builder.WriteString(")")

	return builder.String()
}
