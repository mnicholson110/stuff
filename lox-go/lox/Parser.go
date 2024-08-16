package lox

type ParseError struct {
	token   *Token
	message string
}

type Parser struct {
	tokens  []Token
	current int
}

func NewParser(tokens []Token) *Parser {
	return &Parser{tokens, 0}
}

func NewParseError(token *Token, message string) *ParseError {
	return &ParseError{token, message}
}

func (p *Parser) expression() (Expr, *ParseError) {
	return p.equality()
}

func (p *Parser) equality() (Expr, *ParseError) {
	expr, err := p.comparison()
	if err != nil {
		return nil, err
	}

	for p.match(BANG_EQUAL, EQUAL_EQUAL) {
		operator := p.previous()
		right, err := p.comparison()
		if err != nil {
			return nil, err
		}
		expr = NewBinary(expr, operator, right)
	}

	return expr, nil
}

func (p *Parser) comparison() (Expr, *ParseError) {
	expr, err := p.term()
	if err != nil {
		return nil, err
	}

	for p.match(GREATER, GREATER_EQUAL, LESS, LESS_EQUAL) {
		operator := p.previous()
		right, err := p.term()
		if err != nil {
			return nil, err
		}
		expr = NewBinary(expr, operator, right)
	}

	return expr, nil
}

func (p *Parser) term() (Expr, *ParseError) {
	expr, err := p.factor()
	if err != nil {
		return nil, err
	}

	for p.match(MINUS, PLUS) {
		operator := p.previous()
		right, err := p.factor()
		if err != nil {
			return nil, err
		}
		expr = NewBinary(expr, operator, right)
	}

	return expr, nil
}

func (p *Parser) factor() (Expr, *ParseError) {
	expr, err := p.unary()
	if err != nil {
		return nil, err
	}

	for p.match(SLASH, STAR) {
		operator := p.previous()
		right, err := p.unary()
		if err != nil {
			return nil, err
		}
		expr = NewBinary(expr, operator, right)
	}

	return expr, nil
}

func (p *Parser) unary() (Expr, *ParseError) {
	if p.match(BANG, MINUS) {
		operator := p.previous()
		right, err := p.unary()
		if err != nil {
			return nil, err
		}
		return NewUnary(operator, right), nil
	}

	a, err := p.primary()
	if err != nil {
		return nil, err
	}
	return a, nil
}

func (p *Parser) primary() (Expr, *ParseError) {
	if p.match(FALSE) {
		return NewLiteral(false), nil
	}

	if p.match(TRUE) {
		return NewLiteral(true), nil
	}

	if p.match(NIL) {
		return NewLiteral(nil), nil
	}

	if p.match(NUMBER, STRING) {
		return NewLiteral(p.previous().Literal), nil
	}

	if p.match(LEFT_PAREN) {
		expr, err := p.expression()
		if err != nil {
			return nil, err
		}
		if _, err := p.consume(RIGHT_PAREN, "Expect ')' after expression."); err != nil {
			return nil, err
		}
		return NewGrouping(expr), nil
	}
	err := NewParseError(p.peek(), "Expect expression.")
	err.error()
	return nil, err
}

func (p *Parser) consume(t TokenType, message string) (*Token, *ParseError) {
	if p.check(t) {
		return p.advance(), nil
	}

	err := NewParseError(p.peek(), message)
	err.error()
	return nil, err
}

func (p *Parser) match(types ...TokenType) bool {
	for _, t := range types {
		if p.check(t) {
			p.advance()
			return true
		}
	}

	return false
}

func (p *Parser) check(t TokenType) bool {
	if p.isAtEnd() {
		return false
	}

	return p.peek().Type == t
}

func (p *Parser) advance() *Token {
	if !p.isAtEnd() {
		p.current++
	}

	return p.previous()
}

func (p *Parser) isAtEnd() bool {
	return p.peek().Type == EOF
}

func (p *Parser) peek() *Token {
	return &p.tokens[p.current]
}

func (p *Parser) previous() *Token {
	return &p.tokens[p.current-1]
}

func (p *ParseError) error() {
	ParseErrorHandle(p.token, p.message)
}

func (p *Parser) synchronize() {
	p.advance()

	for !p.isAtEnd() {
		if p.previous().Type == SEMICOLON {
			return
		}

		switch p.peek().Type {
		case CLASS, FUN, VAR, FOR, IF, WHILE, PRINT, RETURN:
			return
		}

		p.advance()
	}
}

func (p *Parser) Parse() (Expr, *ParseError) {
	return p.expression()
}
