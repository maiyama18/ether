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
)

func precedence(t token.Token) Precedence {
	switch t.Type {
	case token.PLUS, token.MINUS:
		return ADDITION
	case token.ASTER, token.SLASH:
		return MULTIPLICATION
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
	if err := p.expectToken(token.SEMICOLON); err != nil {
		return nil, err
	}

	return &ast.VarStatement{Identifier: identifier, Expression: expression}, nil
}

func (p *Parser) parseReturnStatement() (*ast.ReturnStatement, error) {
	p.consumeToken()

	expression, err := p.parseExpression(LOWEST)
	if err != nil {
		return nil, err
	}
	if err := p.expectToken(token.SEMICOLON); err != nil {
		return nil, err
	}

	return &ast.ReturnStatement{Expression: expression}, nil
}

func (p *Parser) parseExpressionStatement() (*ast.ExpressionStatement, error) {
	expression, err := p.parseExpression(LOWEST)
	if err != nil {
		return nil, err
	}
	if err := p.expectToken(token.SEMICOLON); err != nil {
		return nil, err
	}

	return &ast.ExpressionStatement{Expression: expression}, nil
}

func (p *Parser) parseExpression(precedence Precedence) (ast.Expression, error) {
	var left ast.Expression
	var err error
	switch p.currentToken.Type {
	case token.INTEGER:
		left, err = p.parseInteger()
	case token.IDENT:
		left, err = p.parseIdentifier()
	case token.MINUS:
		left, err = p.parsePrefixExpression()
	case token.LPAREN:
		left, err = p.parseGroupedExpression()
	default:
		return nil, &ParserError{line: p.currentToken.Line, msg: fmt.Sprintf("unable to parse token %+v\n", p.currentToken)}
	}
	if err != nil {
		return nil, err
	}

	for precedence < p.peekPrecedence() {
		p.consumeToken()
		left, err = p.parseInfixExpression(left)
		if err != nil {
			return nil, err
		}
	}

	return left, nil
}

func (p *Parser) parseInteger() (*ast.IntegerLiteral, error) {
	v, err := strconv.Atoi(p.currentToken.Literal)
	if err != nil {
		return nil, &ParserError{line: p.currentToken.Line, msg: err.Error()}
	}
	return &ast.IntegerLiteral{Value: v}, nil
}

func (p *Parser) parseIdentifier() (*ast.Identifier, error) {
	return &ast.Identifier{Name: p.currentToken.Literal}, nil
}

func (p *Parser) parsePrefixExpression() (*ast.PrefixExpression, error) {
	operator := p.currentToken.Literal
	p.consumeToken()
	right, err := p.parseExpression(PREFIX)
	if err != nil {
		return nil, err
	}
	return &ast.PrefixExpression{Operator: operator, Right: right}, nil
}

func (p *Parser) parseInfixExpression(left ast.Expression) (*ast.InfixExpression, error) {
	precedence := p.currentPrecedence()
	operator := p.currentToken.Literal
	p.consumeToken()
	right, err := p.parseExpression(precedence)
	if err != nil {
		return nil, err
	}
	return &ast.InfixExpression{Operator: operator, Left: left, Right: right}, nil
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
