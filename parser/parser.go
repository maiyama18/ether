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

func (p *Parser) ParseProgram() *ast.Program {
	statements := make([]ast.Statement, 0)

	for p.currentToken.Type != token.EOF {
		statements = append(statements, p.parseStatement())
		p.consumeToken()
	}

	return &ast.Program{Statements: statements}
}

func (p *Parser) consumeToken() {
	p.currentToken = p.peekToken
	p.peekToken = p.lexer.NextToken()
}

func (p *Parser) expectToken(tokenType token.Type) {
	if p.peekToken.Type != tokenType {
		p.addParserError(
			p.peekToken.Line,
			fmt.Sprintf("expectToken failed.\nwant=%T\ngot=%v (%+v)\n", tokenType, p.peekToken.Type, p.peekToken),
		)
	}
	p.consumeToken()
}

func (p *Parser) currentPrecedence() Precedence {
	return precedence(p.currentToken)
}

func (p *Parser) peekPrecedence() Precedence {
	return precedence(p.peekToken)
}

func (p *Parser) addParserError(line int, msg string) {
	p.errors = append(p.errors, &ParserError{line: line, msg: msg})
}

func (p *Parser) parseStatement() ast.Statement {
	switch p.currentToken.Type {
	case token.VAR:
		return p.parseVarStatement()
	case token.RETURN:
		return p.parseReturnStatement()
	default:
		return p.parseExpressionStatement()
	}
}

func (p *Parser) parseVarStatement() *ast.VarStatement {
	p.consumeToken()

	identifier := p.parseIdentifier()

	p.expectToken(token.ASSIGN)
	p.consumeToken()

	expression := p.parseExpression(LOWEST)
	p.expectToken(token.SEMICOLON)

	return &ast.VarStatement{Identifier: identifier, Expression: expression}
}

func (p *Parser) parseReturnStatement() *ast.ReturnStatement {
	p.consumeToken()

	expression := p.parseExpression(LOWEST)
	p.expectToken(token.SEMICOLON)

	return &ast.ReturnStatement{Expression: expression}
}

func (p *Parser) parseExpressionStatement() *ast.ExpressionStatement {
	expression := p.parseExpression(LOWEST)
	p.expectToken(token.SEMICOLON)

	return &ast.ExpressionStatement{Expression: expression}
}

func (p *Parser) parseExpression(precedence Precedence) ast.Expression {
	var left ast.Expression
	switch p.currentToken.Type {
	case token.INTEGER:
		left = p.parseInteger()
	case token.IDENT:
		left = p.parseIdentifier()
	case token.MINUS:
		left = p.parsePrefixExpression()
	}

	for precedence < p.peekPrecedence() {
		p.consumeToken()
		left = p.parseInfixExpression(left)
	}

	return left
}

func (p *Parser) parseInteger() *ast.IntegerLiteral {
	v, _ := strconv.Atoi(p.currentToken.Literal)
	return &ast.IntegerLiteral{Value: v}
}

func (p *Parser) parseIdentifier() *ast.Identifier {
	return &ast.Identifier{Name: p.currentToken.Literal}
}

func (p *Parser) parsePrefixExpression() *ast.PrefixExpression {
	operator := p.currentToken.Literal
	p.consumeToken()
	right := p.parseExpression(PREFIX)
	return &ast.PrefixExpression{Operator: operator, Right: right}
}

func (p *Parser) parseInfixExpression(left ast.Expression) *ast.InfixExpression {
	precedence := p.currentPrecedence()
	operator := p.currentToken.Literal
	p.consumeToken()
	right := p.parseExpression(precedence)
	return &ast.InfixExpression{Operator: operator, Left: left, Right: right}
}
