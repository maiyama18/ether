package parser

import (
	"fmt"
	"github.com/muiscript/ether/ast"
	"github.com/muiscript/ether/lexer"
	"github.com/muiscript/ether/token"
)

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

func (p *Parser) addParserError(line int, msg string) {
	p.errors = append(p.errors, &ParserError{line: line, msg: msg})
}

func (p *Parser) parseStatement() ast.Statement {
	switch p.currentToken.Type {
	case token.LET:
		return p.parseLetStatement()
	default:
		return nil
	}
}

func (p *Parser) parseLetStatement() *ast.LetStatement {
	tok := p.currentToken
	p.consumeToken()

	identifier := ast.NewIdentifier(p.currentToken)

	p.expectToken(token.ASSIGN)

	// TODO: parse expression
	for p.currentToken.Type != token.SEMICOLON {
		p.consumeToken()
	}

	return ast.NewLetStatement(tok, identifier, nil)
}
