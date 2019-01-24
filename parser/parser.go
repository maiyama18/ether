package parser

import (
	"fmt"
	"github.com/muiscript/ether/ast"
	"github.com/muiscript/ether/lexer"
	"github.com/muiscript/ether/token"
	"strconv"
)

type Precedence int

const (
	LOWEST Precedence = iota
	ADDITION
	MULTIPLICATION
	PREFIX
	CALL
)

func precedence(t token.Token) Precedence {
	switch t.Type {
	case token.PLUS, token.MINUS:
		return ADDITION
	case token.ASTER, token.SLASH:
		return MULTIPLICATION
	case token.LPAREN:
		return PREFIX
	default:
		return LOWEST
	}
}

type Parser struct {
	lexer        *lexer.Lexer
	currentToken token.Token
	peekToken    token.Token
	errors       []*ParserError
}

func New(lexer *lexer.Lexer) *Parser {
	parser := &Parser{lexer: lexer}
	parser.consumeToken()
	parser.consumeToken()

	return parser
}

func (p *Parser) ParseProgram() (*ast.Program, error) {
	statements := make([]ast.Statement, 0)

	for p.currentToken.Type != token.EOF {
		statement, err := p.parseStatement()
		if err != nil {
			return nil, err
		}
		statements = append(statements, statement)
		p.consumeToken()
	}

	return &ast.Program{Statements: statements}, nil
}

func (p *Parser) consumeToken() {
	p.currentToken = p.peekToken
	p.peekToken = p.lexer.NextToken()
}

func (p *Parser) expectToken(tokenType token.Type) error {
	if p.peekToken.Type != tokenType {
		return &ParserError{
			line: p.peekToken.Line,
			msg:  fmt.Sprintf("unexpected token.\nwant=%v\ngot=%v (%+v)\n", tokenType, p.peekToken.Type, p.peekToken),
		}
	}
	p.consumeToken()
	return nil
}

func (p *Parser) currentPrecedence() Precedence {
	return precedence(p.currentToken)
}

func (p *Parser) peekPrecedence() Precedence {
	return precedence(p.peekToken)
}

func (p *Parser) parseStatement() (ast.Statement, error) {
	switch p.currentToken.Type {
	case token.VAR:
		return p.parseVarStatement()
	case token.RETURN:
		return p.parseReturnStatement()
	default:
		return p.parseExpressionStatement()
	}
}

func (p *Parser) parseVarStatement() (*ast.VarStatement, error) {
	line := p.currentToken.Line
	p.consumeToken()

	identifier, err := p.parseIdentifier()
	if err != nil {
		return nil, err
	}

	if err := p.expectToken(token.ASSIGN); err != nil {
		return nil, err
	}
	p.consumeToken()

	expression, err := p.parseExpression(LOWEST)
	if err != nil {
		return nil, err
	}
	if p.peekToken.Type == token.SEMICOLON {
        p.consumeToken()
	}

	return ast.NewVarStatement(identifier, expression, line), nil
}

func (p *Parser) parseReturnStatement() (*ast.ReturnStatement, error) {
	line := p.currentToken.Line
	p.consumeToken()

	expression, err := p.parseExpression(LOWEST)
	if err != nil {
		return nil, err
	}
	if p.peekToken.Type == token.SEMICOLON {
		p.consumeToken()
	}

	return ast.NewReturnStatement(expression, line), nil
}

func (p *Parser) parseExpressionStatement() (*ast.ExpressionStatement, error) {
	line := p.currentToken.Line
	expression, err := p.parseExpression(LOWEST)
	if err != nil {
		return nil, err
	}
	if p.peekToken.Type == token.SEMICOLON {
		p.consumeToken()
	}

	return ast.NewExpressionStatement(expression, line), nil
}

func (p *Parser) parseBlockStatement() (*ast.BlockStatement, error) {
	line := p.currentToken.Line
	p.consumeToken()
	statements := make([]ast.Statement, 0)

	for p.currentToken.Type != token.RBRACE {
		statement, err := p.parseStatement()
		if err != nil {
			return nil, err
		}
		statements = append(statements, statement)
		p.consumeToken()
	}

	return ast.NewBlockStatement(statements, line), nil
}

func (p *Parser) parseExpression(precedence Precedence) (ast.Expression, error) {
	var left ast.Expression
	var err error
	switch p.currentToken.Type {
	case token.INTEGER:
		left, err = p.parseIntegerLiteral()
	case token.IDENT:
		left, err = p.parseIdentifier()
	case token.MINUS:
		left, err = p.parsePrefixExpression()
	case token.LPAREN:
		left, err = p.parseGroupedExpression()
	case token.BAR:
		left, err = p.parseFunctionLiteral()
	default:
		return nil, &ParserError{line: p.currentToken.Line, msg: fmt.Sprintf("unable to parse prefix token %+v\n", p.currentToken)}
	}
	if err != nil {
		return nil, err
	}

	for precedence < p.peekPrecedence() {
		p.consumeToken()
		switch p.currentToken.Type {
		case token.LPAREN:
			left, err = p.parseFunctionCall(left)
		default:
			left, err = p.parseInfixExpression(left)
		}
		if err != nil {
			return nil, err
		}
	}

	return left, nil
}

func (p *Parser) parseIntegerLiteral() (*ast.IntegerLiteral, error) {
	line := p.currentToken.Line
	v, err := strconv.Atoi(p.currentToken.Literal)
	if err != nil {
		return nil, &ParserError{line: line, msg: err.Error()}
	}
	return ast.NewIntegerLiteral(v, line), nil
}

func (p *Parser) parseIdentifier() (*ast.Identifier, error) {
	line := p.currentToken.Line
	if p.currentToken.Type != token.IDENT {
		return nil, &ParserError{line: line, msg: fmt.Sprintf("not identifier: %+v", p.currentToken)}
	}
	return ast.NewIdentifier(p.currentToken.Literal, line), nil
}

func (p *Parser) parsePrefixExpression() (*ast.PrefixExpression, error) {
	line := p.currentToken.Line
	operator := p.currentToken.Literal
	p.consumeToken()
	right, err := p.parseExpression(PREFIX)
	if err != nil {
		return nil, err
	}
	return ast.NewPrefixExpression(operator, right, line), nil
}

func (p *Parser) parseGroupedExpression() (ast.Expression, error) {
	p.consumeToken()
	expression, err := p.parseExpression(LOWEST)
	if err != nil {
		return nil, err
	}
	if err := p.expectToken(token.RPAREN); err != nil {
		return nil, err
	}
	return expression, nil
}

func (p *Parser) parseFunctionLiteral() (ast.Expression, error) {
	line := p.currentToken.Line
	expressions, err := p.parseCommaSeparatedExpressions(token.BAR)
	if err != nil {
		return nil, err
	}
	var parameters []*ast.Identifier
	for _, expression := range expressions {
		if parameter, ok := expression.(*ast.Identifier); ok {
			parameters = append(parameters, parameter)
		} else {
			return nil, &ParserError{line: line, msg: fmt.Sprintf("unable to parse function parameter: %+v", expression)}
		}
	}

	if err := p.expectToken(token.LBRACE); err != nil {
		return nil, err
	}
	body, err := p.parseBlockStatement()
	if err != nil {
		return nil, err
	}

	return ast.NewFunctionLiteral(parameters, body, line), nil
}

func (p *Parser) parseInfixExpression(left ast.Expression) (*ast.InfixExpression, error) {
	line := p.currentToken.Line
	precedence := p.currentPrecedence()
	operator := p.currentToken.Literal
	p.consumeToken()
	right, err := p.parseExpression(precedence)
	if err != nil {
		return nil, err
	}
	return ast.NewInfixExpression(operator, left, right, line), nil
}

func (p *Parser) parseFunctionCall(left ast.Expression) (*ast.FunctionCall, error) {
	line := p.currentToken.Line
	arguments, err := p.parseCommaSeparatedExpressions(token.RPAREN)
	if err != nil {
		return nil, err
	}
	return ast.NewFunctionCall(left, arguments, line), nil
}

func (p *Parser) parseCommaSeparatedExpressions(endTokenType token.Type) ([]ast.Expression, error) {
	p.consumeToken()
	if p.currentToken.Type == endTokenType {
		return []ast.Expression{}, nil
	}

	first, err := p.parseExpression(LOWEST)
	if err != nil {
		return nil, err
	}
	expressions := []ast.Expression{first}

	for p.peekToken.Type != endTokenType {
		if err := p.expectToken(token.COMMA); err != nil {
			return nil, err
		}
		p.consumeToken()

		expression, err := p.parseExpression(LOWEST)
		if err != nil {
			return nil, err
		}
		expressions = append(expressions, expression)
	}
	if err := p.expectToken(endTokenType); err != nil {
		return nil, err
	}

	return expressions, nil
}

