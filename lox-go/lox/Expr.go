package lox

import (
	"strconv"
)

type Expr interface {
	AstPrint() string
}

type Binary struct {
	left     Expr
	operator *Token
	right    Expr
}

type Grouping struct {
	expression Expr
}

type Literal struct {
	value interface{}
}

type Unary struct {
	operator *Token
	right    Expr
}

func (b *Binary) AstPrint() string {
	return AstPrinter(b.operator.Lexeme, b.left, b.right)
}

func (g *Grouping) AstPrint() string {
	return AstPrinter("group", g.expression)
}

func (l *Literal) AstPrint() string {
	if l == nil {
		return "nil"
	}
	switch l.value.(type) {
	case string:
		return l.value.(string)
	case float64:
		v, ok := l.value.(float64)
		if ok {
			return strconv.FormatFloat(v, 'f', -1, 64)
		}
	case bool:
		v, ok := l.value.(bool)
		if ok {
			return strconv.FormatBool(v)
		}
	case int:
		v, ok := l.value.(int)
		if ok {
			return strconv.Itoa(v)
		}
	}
	return "undefined"
}

func (u *Unary) AstPrint() string {
	return AstPrinter(u.operator.Lexeme, u.right)
}

// Constructors
func NewBinary(left Expr, operator *Token, right Expr) *Binary {
	return &Binary{left, operator, right}
}

func NewGrouping(expression Expr) *Grouping {
	return &Grouping{expression}
}

func NewLiteral(value interface{}) *Literal {
	return &Literal{value}
}

func NewUnary(operator *Token, right Expr) *Unary {
	return &Unary{operator, right}
}
