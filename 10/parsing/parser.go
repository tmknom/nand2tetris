package parsing

import (
	"../token"
	"fmt"
	"github.com/pkg/errors"
)

type Parser struct {
	tokens *token.Tokens
}

func NewParser(tokens *token.Tokens) *Parser {
	return &Parser{tokens: tokens}
}

func (p *Parser) advanceToken() *token.Token {
	return p.tokens.Advance()
}

//func (p *Parser) backwardToken() *token.Token {
//	return p.tokens.Backward()
//}

func (p *Parser) readFirstToken() *token.Token {
	return p.tokens.First()
}

func (p *Parser) readSecondToken() *token.Token {
	return p.tokens.Second()
}

func (p *Parser) Parse() (*Class, error) {
	class, err := p.parseClass()
	if err != nil {
		return nil, errors.WithMessage(err, p.tokens.DebugForError())
	}

	return class, nil
}

// 'class' className '{' classVarDec* subroutineDec* '}'
// class Main { ... }
func (p *Parser) parseClass() (*Class, error) {
	class := NewClass()

	keyword := p.advanceToken()
	if err := class.CheckKeyword(keyword); err != nil {
		return nil, err
	}

	className := p.advanceToken()
	if err := class.SetClassName(className); err != nil {
		return nil, err
	}

	openingCurlyBracket := p.advanceToken()
	if err := ConstOpeningCurlyBracket.Check(openingCurlyBracket); err != nil {
		return nil, err
	}

	classVarDecs, err := p.parseClassVarDecs()
	if err != nil {
		return nil, err
	}
	class.SetClassVarDecs(classVarDecs)

	subroutineDecs, err := p.parseSubroutineDecs()
	if err != nil {
		return nil, err
	}
	class.SetSubroutineDecs(subroutineDecs)

	closingCurlyBracket := p.advanceToken()
	if err := ConstClosingCurlyBracket.Check(closingCurlyBracket); err != nil {
		return nil, err
	}

	return class, nil
}

// ('static' | 'field') varType varName (',' varName) ';'
// field int x, y;
func (p *Parser) parseClassVarDecs() (*ClassVarDecs, error) {
	classVarDecs := NewClassVarDecs()

	for classVarDecs.HasClassVarDec(p.readFirstToken()) {
		classVarDec := NewClassVarDec()

		keyword := p.advanceToken()
		if err := classVarDec.SetKeyword(keyword); err != nil {
			return nil, err
		}

		varType := p.advanceToken()
		if err := classVarDec.SetVarType(varType); err != nil {
			return nil, err
		}

		varName := p.advanceToken()
		if err := classVarDec.SetFirstVarName(varName); err != nil {
			return nil, err
		}

		for ConstComma.IsCheck(p.readFirstToken()) {
			comma := p.advanceToken()
			varName := p.advanceToken()
			if err := classVarDec.AddCommaAndVarName(comma, varName); err != nil {
				return nil, err
			}
		}

		semicolon := p.advanceToken()
		if err := ConstSemicolon.Check(semicolon); err != nil {
			return nil, err
		}

		// パースに成功したら要素に追加
		classVarDecs.Add(classVarDec)
	}

	return classVarDecs, nil
}

// ('constructor' | 'function' | 'method') ('void' | varType) subroutineName '(' parameterList ')' subroutineBody
// constructor Square new(int x, int y) { ... }
func (p *Parser) parseSubroutineDecs() (*SubroutineDecs, error) {
	subroutineDecs := NewSubroutineDecs()
	for subroutineDecs.hasSubroutineDec(p.readFirstToken()) {
		keyword := NewKeyword(p.advanceToken())
		subroutineDec := NewSubroutineDec(keyword)

		subroutineType := p.advanceToken()
		if err := subroutineDec.SetSubroutineType(subroutineType); err != nil {
			return nil, err
		}

		subroutineName := p.advanceToken()
		if err := subroutineDec.SetSubroutineName(subroutineName); err != nil {
			return nil, err
		}

		openingRoundBracket := p.advanceToken()
		if err := ConstOpeningRoundBracket.Check(openingRoundBracket); err != nil {
			return nil, err
		}

		// パラメータリストの追加
		parameterList, err := p.parseParameterList()
		if err != nil {
			return nil, err
		}
		subroutineDec.SetParameterList(parameterList)

		closingRoundBracket := p.advanceToken()
		if err := ConstClosingRoundBracket.Check(closingRoundBracket); err != nil {
			return nil, err
		}

		subroutineBody, err := p.parseSubroutineBody()
		if err != nil {
			return nil, err
		}
		subroutineDec.SetSubroutineBody(subroutineBody)

		// パースに成功したら要素に追加
		subroutineDecs.Add(subroutineDec)
	}

	return subroutineDecs, nil
}

// ((varType varName) (',' varType varName)*)?
// int Ax, int Ay
func (p *Parser) parseParameterList() (*ParameterList, error) {
	// パラメータがひとつも定義されていない場合は即終了
	parameterList := NewParameterList()
	if !NewVarType(p.readFirstToken()).IsCheck() {
		return parameterList, nil
	}

	// パラメータ1つめのみカンマがないのでループに入る前に処理する
	varType := p.advanceToken()
	varName := p.advanceToken()
	if err := parameterList.Add(varType, varName); err != nil {
		return nil, err
	}

	// パラメータ2つめ以降はカンマが見つかった場合のみ処理する
	for ConstComma.IsCheck(p.readFirstToken()) {
		p.advanceToken() // カンマを飛ばす
		varType := p.advanceToken()
		varName := p.advanceToken()
		if err := parameterList.Add(varType, varName); err != nil {
			return nil, err
		}
	}
	return parameterList, nil
}

// '{' varDec* statements '}'
func (p *Parser) parseSubroutineBody() (*SubroutineBody, error) {
	subroutineBody := NewSubroutineBody()

	openingCurlyBracket := p.advanceToken()
	if err := ConstOpeningCurlyBracket.Check(openingCurlyBracket); err != nil {
		return nil, err
	}

	// varDecのパース
	for subroutineBody.IsVarDecKeyword(p.readFirstToken()) {
		varDec, err := p.parseVarDec()
		if err != nil {
			return nil, err
		}
		subroutineBody.AddVarDec(varDec)
	}

	// statementsのパース
	for subroutineBody.IsStatementKeyword(p.readFirstToken()) {
		statement, err := p.parseStatement()
		if err != nil {
			return nil, err
		}
		subroutineBody.AddStatement(statement)
	}

	closingCurlyBracket := p.advanceToken()
	if err := ConstClosingCurlyBracket.Check(closingCurlyBracket); err != nil {
		return nil, err
	}

	return subroutineBody, nil
}

// 'var' varType varName (',' varName) ';'
func (p *Parser) parseVarDec() (*VarDec, error) {
	varDec := NewVarDec()

	keyword := p.advanceToken()
	if err := varDec.CheckKeyword(keyword); err != nil {
		return nil, err
	}

	varType := p.advanceToken()
	if err := varDec.SetVarType(varType); err != nil {
		return nil, err
	}

	varName := p.advanceToken()
	if err := varDec.SetFirstVarName(varName); err != nil {
		return nil, err
	}

	for ConstComma.IsCheck(p.readFirstToken()) {
		comma := p.advanceToken()
		varName := p.advanceToken()
		if err := varDec.AddCommaAndVarName(comma, varName); err != nil {
			return nil, err
		}
	}

	semicolon := p.advanceToken()
	if err := ConstSemicolon.Check(semicolon); err != nil {
		return nil, err
	}

	return varDec, nil
}

func (p *Parser) parseStatement() (Statement, error) {
	keyword := p.advanceToken()
	switch keyword.Value {
	case "let":
		return p.parseNotImplementedStatement()
	case "if":
		return p.parseNotImplementedStatement()
	case "while":
		return p.parseNotImplementedStatement()
	case "do":
		return p.parseNotImplementedStatement()
	case "return":
		return p.parseReturnStatement()
	default:
		message := fmt.Sprintf("Invalid Statement: got = %s", keyword.Debug())
		return nil, errors.New(message)
	}
}

// (let) varName ('[' expression ']')? '=' expression ';'
func (p *Parser) parseLetStatement() (Statement, error) {
	letStatement := NewLetStatement()

	varName := p.advanceToken()
	if err := letStatement.SetVarName(varName); err != nil {
		return nil, err
	}

	if ConstOpeningSquareBracket.IsCheck(p.readFirstToken()) {
		p.advanceToken() // '[' を飛ばす

		arrayExpression := p.advanceToken()
		if err := letStatement.SetArrayIndex(arrayExpression); err != nil {
			return nil, err
		}

		closingSquareBracket := p.advanceToken()
		if err := ConstClosingSquareBracket.Check(closingSquareBracket); err != nil {
			return nil, err
		}
	}

	equal := p.advanceToken()
	if err := ConstEqual.Check(equal); err != nil {
		return nil, err
	}

	expression := p.advanceToken()
	if err := letStatement.SetExpression(expression); err != nil {
		return nil, err
	}

	semicolon := p.advanceToken()
	if err := ConstSemicolon.Check(semicolon); err != nil {
		return nil, err
	}

	return letStatement, nil
}

// (do) subroutineCall ';'
func (p *Parser) parseDoStatement() (Statement, error) {
	doStatement := NewDoStatement()

	subroutineCall, err := p.parseSubroutineCall()
	if err != nil {
		return nil, err
	}
	doStatement.SetSubroutineCall(subroutineCall)

	semicolon := p.advanceToken()
	if err := ConstSemicolon.Check(semicolon); err != nil {
		return nil, err
	}

	return doStatement, nil
}

// (return) expression? ';'
func (p *Parser) parseReturnStatement() (Statement, error) {
	returnStatement := NewReturnStatement()

	if !ConstSemicolon.IsCheck(p.readFirstToken()) {
		expression := p.advanceToken()
		if err := returnStatement.SetExpression(expression); err != nil {
			return nil, err
		}
	}

	semicolon := p.advanceToken()
	if err := ConstSemicolon.Check(semicolon); err != nil {
		return nil, err
	}

	return returnStatement, nil
}

// subroutineName '(' expressionList ')'
// (className | varName) '.' subroutineName '(' expressionList ')'
func (p *Parser) parseSubroutineCall() (*SubroutineCall, error) {
	subroutineCall := NewSubroutineCall()

	subroutineCallName, err := p.parseSubroutineCallName()
	if err != nil {
		return nil, err
	}
	subroutineCall.SetSubroutineCallName(subroutineCallName)

	openingRoundBracket := p.advanceToken()
	if err := ConstOpeningRoundBracket.Check(openingRoundBracket); err != nil {
		return nil, err
	}

	expressionList, err := p.parseExpressionList()
	if err != nil {
		return nil, err
	}
	subroutineCall.SetExpressionList(expressionList)

	closingRoundBracket := p.advanceToken()
	if err := ConstClosingRoundBracket.Check(closingRoundBracket); err != nil {
		return nil, err
	}

	return subroutineCall, nil
}

// subroutineName
// (className | varName) '.' subroutineName
func (p *Parser) parseSubroutineCallName() (*SubroutineCallName, error) {
	subroutineCallName := NewSubroutineCallName()
	name := p.advanceToken()

	if ConstPeriod.IsCheck(p.readFirstToken()) {
		if err := subroutineCallName.SetCallerName(name); err != nil {
			return nil, err
		}

		period := p.advanceToken()
		if err := ConstPeriod.Check(period); err != nil {
			return nil, err
		}

		name = p.advanceToken()
	}

	if err := subroutineCallName.SetSubroutineName(name); err != nil {
		return nil, err
	}

	return subroutineCallName, nil
}

func (p *Parser) parseExpressionList() (*ExpressionList, error) {
	// 式がひとつも定義されていない場合は即終了
	expressionList := NewExpressionList()
	if !NewExpression(p.readFirstToken()).IsCheck() {
		return expressionList, nil
	}

	// 1つめのみカンマがないのでループに入る前に処理する
	expression := p.advanceToken()
	if err := expressionList.AddExpression(expression); err != nil {
		return nil, err
	}

	// 2つめ以降はカンマが見つかった場合のみ処理する
	for ConstComma.IsCheck(p.readFirstToken()) {
		p.advanceToken() // カンマを飛ばす
		expression := p.advanceToken()
		if err := expressionList.AddExpression(expression); err != nil {
			return nil, err
		}
	}
	return expressionList, nil
}

// integerConstant | stringConstant | keywordConstant |
// varName | subroutineCall | varName '[' expression ']' |
// '(' expression ')' | unaryOp term
func (p *Parser) parseTerm() (Term, error) {
	term := p.readFirstToken()
	switch term.TokenType {
	case token.TokenIntConst:
		return p.parseIntegerConstant()
	case token.TokenStringConst:
		return p.parseStringConstant()
	case token.TokenKeyword:
		return p.parseKeywordConstant()
	case token.TokenIdentifier:
		return p.parseIdentifierTerm()
	case token.TokenSymbol:
		return p.parseSymbolTerm()
	default:
		message := fmt.Sprintf("error parseTerm: got = %s", term.Debug())
		return nil, errors.New(message)
	}
}

// integerConstant
func (p *Parser) parseIntegerConstant() (Term, error) {
	integerConstant := NewIntegerConstant(p.advanceToken())
	if err := integerConstant.Check(); err != nil {
		return nil, err
	}
	return integerConstant, nil
}

// stringConstant
func (p *Parser) parseStringConstant() (Term, error) {
	stringConstant := NewStringConstant(p.advanceToken())
	if err := stringConstant.Check(); err != nil {
		return nil, err
	}
	return stringConstant, nil
}

// keywordConstant
func (p *Parser) parseKeywordConstant() (Term, error) {
	return ConstKeywordConstantFactory.Create(p.advanceToken())
}

// varName | subroutineCall | varName '[' expression ']'
func (p *Parser) parseIdentifierTerm() (Term, error) {
	second := p.readSecondToken()

	switch second.Value {
	case ConstOpeningRoundBracket.Value, ConstPeriod.Value:
		return p.parseSubroutineCall()
	case ConstOpeningSquareBracket.Value:
		return p.parseArray()
	default:
		varName := p.advanceToken()
		return NewVarNameOrError(varName)
	}
}

// '(' expression ')' | unaryOp term
func (p *Parser) parseSymbolTerm() (Term, error) {
	op := p.readFirstToken()

	switch op.Value {
	case ConstMinus.Value, ConstTilde.Value:
		return p.parseUnaryOpTerm()
	default:
		// TODO '(' expression ')'
		message := fmt.Sprintf("not implemented BinaryOp: got = %s", op.Debug())
		return nil, errors.New(message)
	}
}

// unaryOp term
func (p *Parser) parseUnaryOpTerm() (*UnaryOpTerm, error) {
	unary, err := ConstUnaryOpFactory.Create(p.advanceToken())
	if err != nil {
		return nil, err
	}
	unaryOpTerm := NewUnaryOpTerm(unary)

	term, err := p.parseTerm()
	if err != nil {
		return nil, err
	}
	unaryOpTerm.SetTerm(term)

	return unaryOpTerm, nil
}

// varName '[' expression ']'
func (p *Parser) parseArray() (*Array, error) {
	array, err := NewArrayOrError(p.advanceToken())
	if err != nil {
		return nil, err
	}

	openingSquareBracket := p.advanceToken()
	if err := ConstOpeningSquareBracket.Check(openingSquareBracket); err != nil {
		return nil, err
	}

	expression := p.advanceToken()
	if err := array.SetExpression(expression); err != nil {
		return nil, err
	}

	closingSquareBracket := p.advanceToken()
	if err := ConstClosingSquareBracket.Check(closingSquareBracket); err != nil {
		return nil, err
	}

	return array, nil
}

func (p *Parser) parseNotImplementedStatement() (Statement, error) {
	//p.readFirstToken()
	//fmt.Println(p.tokens.Debug())
	//fmt.Println(p.readFirstToken().Debug())

	// TODO parseStatementの実装が完了するまで辻褄をあわせる
	for {
		switch p.readFirstToken().Value {
		case "let", "if", "while", "do", "return":
			return p.parseStatement()
		default:
			p.advanceToken()
		}
	}
}
