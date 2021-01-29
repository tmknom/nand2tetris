package parsing

import (
	"../token"
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
	p.tokens = p.tokens.SubList()
	return p.tokens.First()
}

func (p *Parser) Parse() (*Class, error) {
	return p.parseClass()
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

func (p *Parser) parseStatement() (*Statement, error) {
	statement := NewStatement()

	for {
		// TODO とりあえず実装が完了するまで「return;」まで読み込んで終了する
		t := p.advanceToken()
		if t.Value == "return" {
			end := p.advanceToken() // セミコロンをスキップ
			if end.Value != ";" {
				p.advanceToken() // リターンで値を返す場合にはさらにもう一つトークンをスキップする
			}
			break
		}
	}

	//p.readFirstToken()
	//fmt.Println(p.tokens.Debug())
	//fmt.Println(p.readFirstToken().Debug())

	return statement, nil
}
