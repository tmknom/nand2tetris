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
	statements, err := p.parseStatements()
	if err != nil {
		return nil, err
	}
	subroutineBody.SetStatements(statements)

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

func (p *Parser) parseStatements() (*Statements, error) {
	statements := NewStatements()

	for statements.IsStatementKeyword(p.readFirstToken()) {
		statement, err := p.parseStatement()
		if err != nil {
			return nil, err
		}
		statements.AddStatement(statement)
	}

	return statements, nil
}

func (p *Parser) parseStatement() (Statement, error) {
	keyword := p.readFirstToken()
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

	keyword := p.advanceToken()
	if err := letStatement.CheckKeyword(keyword); err != nil {
		return nil, err
	}

	second := p.readSecondToken()
	if ConstOpeningSquareBracket.IsCheck(second) {
		array, err := p.parseArray()
		if err != nil {
			return nil, err
		}
		letStatement.SetArray(array)
	} else {
		varName := p.advanceToken()
		if err := letStatement.SetVarName(varName); err != nil {
			return nil, err
		}
	}

	equal := p.advanceToken()
	if err := ConstEqual.Check(equal); err != nil {
		return nil, err
	}

	expression, err := p.parseExpression()
	if err != nil {
		return nil, err
	}
	letStatement.SetExpression(expression)

	semicolon := p.advanceToken()
	if err := ConstSemicolon.Check(semicolon); err != nil {
		return nil, err
	}

	return letStatement, nil
}

// if '(' expression ')' '{' statements '}' ( else '{' statements '}' )?
func (p *Parser) parseIfStatement() (Statement, error) {
	ifStatement := NewIfStatement()

	keyword := p.advanceToken()
	if err := ifStatement.CheckKeyword(keyword); err != nil {
		return nil, err
	}

	openingRoundBracket := p.advanceToken()
	if err := ConstOpeningRoundBracket.Check(openingRoundBracket); err != nil {
		return nil, err
	}

	expression, err := p.parseExpression()
	if err != nil {
		return nil, err
	}
	ifStatement.SetExpression(expression)

	closingRoundBracket := p.advanceToken()
	if err := ConstClosingRoundBracket.Check(closingRoundBracket); err != nil {
		return nil, err
	}

	openingCurlyBracket := p.advanceToken()
	if err := ConstOpeningCurlyBracket.Check(openingCurlyBracket); err != nil {
		return nil, err
	}

	statements, err := p.parseStatements()
	if err != nil {
		return nil, err
	}
	ifStatement.SetStatements(statements)

	closingCurlyBracket := p.advanceToken()
	if err := ConstClosingCurlyBracket.Check(closingCurlyBracket); err != nil {
		return nil, err
	}

	// else句が存在するかチェックする
	if NewKeyword(p.readFirstToken()).Check("else") == nil {
		elseBlock, err := p.parseElseBlock()
		if err != nil {
			return nil, err
		}
		ifStatement.SetElseBlock(elseBlock)
	}

	return ifStatement, nil
}

// while '(' expression ')' '{' statements '}'
func (p *Parser) parseWhileStatement() (Statement, error) {
	whileStatement := NewWhileStatement()

	keyword := p.advanceToken()
	if err := whileStatement.CheckKeyword(keyword); err != nil {
		return nil, err
	}

	openingRoundBracket := p.advanceToken()
	if err := ConstOpeningRoundBracket.Check(openingRoundBracket); err != nil {
		return nil, err
	}

	expression, err := p.parseExpression()
	if err != nil {
		return nil, err
	}
	whileStatement.SetExpression(expression)

	closingRoundBracket := p.advanceToken()
	if err := ConstClosingRoundBracket.Check(closingRoundBracket); err != nil {
		return nil, err
	}

	openingCurlyBracket := p.advanceToken()
	if err := ConstOpeningCurlyBracket.Check(openingCurlyBracket); err != nil {
		return nil, err
	}

	statements, err := p.parseStatements()
	if err != nil {
		return nil, err
	}
	whileStatement.SetStatements(statements)

	closingCurlyBracket := p.advanceToken()
	if err := ConstClosingCurlyBracket.Check(closingCurlyBracket); err != nil {
		return nil, err
	}

	return whileStatement, nil
}

// (do) subroutineCall ';'
func (p *Parser) parseDoStatement() (Statement, error) {
	doStatement := NewDoStatement()

	keyword := p.advanceToken()
	if err := doStatement.CheckKeyword(keyword); err != nil {
		return nil, err
	}

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

// return expression? ';'
func (p *Parser) parseReturnStatement() (Statement, error) {
	returnStatement := NewReturnStatement()

	keyword := p.advanceToken()
	if err := returnStatement.CheckKeyword(keyword); err != nil {
		return nil, err
	}

	if !ConstSemicolon.IsCheck(p.readFirstToken()) {
		expression, err := p.parseExpression()
		if err != nil {
			return nil, err
		}
		returnStatement.SetExpression(expression)
	}

	semicolon := p.advanceToken()
	if err := ConstSemicolon.Check(semicolon); err != nil {
		return nil, err
	}

	return returnStatement, nil
}

// else '{' statements '}'
func (p *Parser) parseElseBlock() (*ElseBlock, error) {
	elseBlock := NewElseBlock()

	keyword := p.advanceToken()
	if err := elseBlock.CheckElseKeyword(keyword); err != nil {
		return nil, err
	}

	openingCurlyBracket := p.advanceToken()
	if err := ConstOpeningCurlyBracket.Check(openingCurlyBracket); err != nil {
		return nil, err
	}

	statements, err := p.parseStatements()
	if err != nil {
		return nil, err
	}
	elseBlock.SetStatements(statements)

	closingCurlyBracket := p.advanceToken()
	if err := ConstClosingCurlyBracket.Check(closingCurlyBracket); err != nil {
		return nil, err
	}

	return elseBlock, nil
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
	if ConstClosingRoundBracket.IsCheck(p.readFirstToken()) {
		return expressionList, nil
	}

	// 1つめのみカンマがないのでループに入る前に処理する
	expression, err := p.parseExpression()
	if err != nil {
		return nil, err
	}
	expressionList.AddExpression(expression)

	// 2つめ以降はカンマが見つかった場合のみ処理する
	for ConstComma.IsCheck(p.readFirstToken()) {
		p.advanceToken() // カンマを飛ばす
		expression, err := p.parseExpression()
		if err != nil {
			return nil, err
		}
		expressionList.AddExpression(expression)
	}
	return expressionList, nil
}

// TODO (op term)* を実装する
// term (op term)*
func (p *Parser) parseExpression() (*Expression, error) {
	term, err := p.parseTerm()
	if err != nil {
		return nil, err
	}
	return NewExpression(term), nil
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
	case ConstOpeningRoundBracket.Value:
		return p.parseGroupingExpression()
	default:
		message := fmt.Sprintf("error parseSymbolTerm: got = %s", op.Debug())
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

// '(' expression ')'
func (p *Parser) parseGroupingExpression() (*GroupingExpression, error) {
	if err := ConstOpeningRoundBracket.Check(p.advanceToken()); err != nil {
		return nil, err
	}

	expression, err := p.parseExpression()
	if err != nil {
		return nil, err
	}
	groupingExpression := NewGroupingExpression(expression)

	if err := ConstClosingRoundBracket.Check(p.advanceToken()); err != nil {
		return nil, err
	}

	return groupingExpression, nil
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

	expression, err := p.parseExpression()
	if err != nil {
		return nil, err
	}
	array.SetExpression(expression)

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
	p.advanceToken()
	for {
		switch p.readFirstToken().Value {
		case "let", "if", "while", "do", "return":
			return p.parseStatement()
		default:
			p.advanceToken()
		}
	}
}
