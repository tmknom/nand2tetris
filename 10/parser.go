package main

import (
	"fmt"
	"github.com/pkg/errors"
)

type Parser struct {
	tokens *Tokens
}

func NewParser(tokens *Tokens) *Parser {
	return &Parser{tokens: tokens}
}

func (p *Parser) advanceToken() *Token {
	return p.tokens.Advance()
}

func (p *Parser) backwardToken() *Token {
	return p.tokens.Backward()
}

func (p *Parser) readFirstToken() *Token {
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

	openSymbol := p.advanceToken()
	if err := ConstOpeningCurlyBracket.Check(openSymbol); err != nil {
		return nil, err
	}

	// 閉じカッコは後ろから取得
	closeSymbol := p.backwardToken()
	if err := ConstClosingCurlyBracket.Check(closeSymbol); err != nil {
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

	return class, nil
}

type Class struct {
	Keyword        *Keyword
	ClassName      *ClassName
	OpenSymbol     *OpeningCurlyBracket
	CloseSymbol    *ClosingCurlyBracket
	ClassVarDecs   *ClassVarDecs
	SubroutineDecs *SubroutineDecs
}

func NewClass() *Class {
	return &Class{
		Keyword:     NewKeywordByValue("class"),
		OpenSymbol:  ConstOpeningCurlyBracket,
		CloseSymbol: ConstClosingCurlyBracket,
	}
}

func (c *Class) CheckKeyword(token *Token) error {
	return NewKeyword(token).Check(c.Keyword.Value)
}

type ClassName struct {
	*Identifier
}

func NewClassName(token *Token) *ClassName {
	return &ClassName{
		Identifier: NewIdentifier("ClassName", token),
	}
}

func (c *Class) SetClassName(token *Token) error {
	className := NewClassName(token)
	if err := className.Check(); err != nil {
		return err
	}

	c.ClassName = className
	return nil
}

func (c *Class) SetClassVarDecs(classVarDecs *ClassVarDecs) {
	c.ClassVarDecs = classVarDecs
}

func (c *Class) SetSubroutineDecs(subroutineDecs *SubroutineDecs) {
	c.SubroutineDecs = subroutineDecs
}

func (c *Class) ToXML() []string {
	result := []string{}
	result = append(result, "<class>")
	result = append(result, c.Keyword.ToXML())
	result = append(result, c.ClassName.ToXML())
	result = append(result, c.OpenSymbol.ToXML())
	result = append(result, c.ClassVarDecs.ToXML()...)
	result = append(result, c.SubroutineDecs.ToXML()...)
	result = append(result, c.CloseSymbol.ToXML())
	result = append(result, "</class>")
	return result
}

func (c *Class) debug() string {
	return fmt.Sprintf("&Class{\n  Keyword: %s,\n  ClassName: %s},\n  OpenSymbol: %s,\n  CloseSymbol: %s,\n  ClassVarDec: &Tokens{...},\n  SubroutineDec: &Tokens{...}\n}",
		c.Keyword.debug(), c.ClassName.debug(), c.OpenSymbol.debug(), c.CloseSymbol.debug())
}

// ('static' | 'field') varType varName (',' varName) ';'
// field int x, y;
func (p *Parser) parseClassVarDecs() (*ClassVarDecs, error) {
	classVarDecs := NewClassVarDecs()

	for classVarDecs.hasClassVarDec(p.readFirstToken()) {
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

		endSymbol := p.advanceToken()
		if err := ConstSemicolon.Check(endSymbol); err != nil {
			return nil, err
		}

		// パースに成功したら要素に追加
		classVarDecs.Add(classVarDec)
	}

	return classVarDecs, nil
}

type ClassVarDecs struct {
	Items []*ClassVarDec
}

func NewClassVarDecs() *ClassVarDecs {
	return &ClassVarDecs{
		Items: []*ClassVarDec{},
	}
}

func (c *ClassVarDecs) Add(item *ClassVarDec) {
	c.Items = append(c.Items, item)
}

func (c *ClassVarDecs) ToXML() []string {
	result := []string{}
	for _, item := range c.Items {
		result = append(result, item.ToXML()...)
	}
	return result
}

func (c *ClassVarDecs) hasClassVarDec(token *Token) bool {
	if token == nil {
		return false
	}
	return token.Value == "static" || token.Value == "field"
}

type ClassVarDec struct {
	*Keyword
	*VarType
	*VarNames
	EndSymbol *Semicolon
}

func NewClassVarDec() *ClassVarDec {
	return &ClassVarDec{
		VarNames:  NewVarNames(),
		EndSymbol: ConstSemicolon,
	}
}

func (c *ClassVarDec) SetKeyword(token *Token) error {
	if err := c.checkKeyword(token); err != nil {
		return err
	}

	c.Keyword = NewKeywordByValue(token.Value)
	return nil
}

func (c *ClassVarDec) checkKeyword(token *Token) error {
	expected := []string{"static", "field"}
	return token.CheckKeywordValue(expected...)
}

func (c *ClassVarDec) SetVarType(token *Token) error {
	varType := NewVarType(token)
	if err := varType.Check(); err != nil {
		return err
	}

	c.VarType = varType
	return nil
}

func (c *ClassVarDec) ToXML() []string {
	result := []string{}
	result = append(result, "<classVarDec>")
	result = append(result, c.Keyword.ToXML())
	result = append(result, c.VarType.ToXML())
	result = append(result, c.VarNames.ToXML()...)
	result = append(result, c.EndSymbol.ToXML())
	result = append(result, "</classVarDec>")
	return result
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

type SubroutineDecs struct {
	Items []*SubroutineDec
}

func NewSubroutineDecs() *SubroutineDecs {
	return &SubroutineDecs{
		Items: []*SubroutineDec{},
	}
}

func (s *SubroutineDecs) Add(item *SubroutineDec) {
	s.Items = append(s.Items, item)
}

func (s *SubroutineDecs) ToXML() []string {
	result := []string{}
	for _, item := range s.Items {
		result = append(result, item.ToXML()...)
	}
	return result
}

func (s *SubroutineDecs) hasSubroutineDec(token *Token) bool {
	if token == nil {
		return false
	}

	return token.Value == "constructor" || token.Value == "function" || token.Value == "method"
}

type SubroutineDec struct {
	Subroutine *Keyword
	*SubroutineType
	*SubroutineName
	*OpeningRoundBracket
	*ClosingRoundBracket
	*ParameterList
	*SubroutineBody
}

func NewSubroutineDec(subroutine *Keyword) *SubroutineDec {
	return &SubroutineDec{
		Subroutine:          subroutine,
		OpeningRoundBracket: ConstOpeningRoundBracket,
		ClosingRoundBracket: ConstClosingRoundBracket,
	}
}

func (s *SubroutineDec) SetSubroutineType(token *Token) error {
	subroutineType := NewSubroutineType(token)
	if err := subroutineType.Check(); err != nil {
		return err
	}

	s.SubroutineType = subroutineType
	return nil
}

func (s *SubroutineDec) SetSubroutineName(token *Token) error {
	subroutineName := NewSubroutineName(token)
	if err := subroutineName.Check(); err != nil {
		return err
	}

	s.SubroutineName = subroutineName
	return nil
}

func (s *SubroutineDec) SetParameterList(parameterList *ParameterList) {
	s.ParameterList = parameterList
}

func (s *SubroutineDec) SetSubroutineBody(subroutineBody *SubroutineBody) {
	s.SubroutineBody = subroutineBody
}

func (s *SubroutineDec) ToXML() []string {
	result := []string{}
	result = append(result, "<subroutineDec>")
	result = append(result, s.Subroutine.ToXML())
	result = append(result, s.SubroutineType.ToXML())
	result = append(result, s.SubroutineName.ToXML())
	result = append(result, s.OpeningRoundBracket.ToXML())
	result = append(result, s.ParameterList.ToXML()...)
	result = append(result, s.ClosingRoundBracket.ToXML())
	result = append(result, s.SubroutineBody.ToXML()...)
	result = append(result, "</subroutineDec>")
	return result
}

type SubroutineType struct {
	*Token
}

func NewSubroutineType(token *Token) *SubroutineType {
	return &SubroutineType{
		Token: token,
	}
}

func (s *SubroutineType) Check() error {
	if err := NewVarType(s.Token).Check(); err == nil {
		return nil
	}

	expected := []string{"void"}
	if err := s.CheckKeywordValue(expected...); err == nil {
		return nil
	}

	message := fmt.Sprintf("SubroutineType: got = %s", s.debug())
	return errors.New(message)
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

type ParameterList struct {
	First              *Parameter
	CommaAndParameters []*CommaAndParameter
}

func NewParameterList() *ParameterList {
	return &ParameterList{
		CommaAndParameters: []*CommaAndParameter{},
	}
}

func (p *ParameterList) Add(varTypeToken *Token, varNameToken *Token) error {
	parameter := NewParameterByToken(varTypeToken, varNameToken)
	if err := parameter.Check(); err != nil {
		return err
	}

	if p.First == nil {
		p.First = parameter
	} else {
		p.CommaAndParameters = append(p.CommaAndParameters, NewCommaAndParameter(parameter))
	}
	return nil
}

func (p *ParameterList) ToXML() []string {
	result := []string{}
	result = append(result, "<parameterList>")

	if p.First != nil {
		result = append(result, p.First.ToXML()...)
	}

	for _, commaAndParameter := range p.CommaAndParameters {
		result = append(result, commaAndParameter.ToXML()...)
	}

	result = append(result, "</parameterList>")
	return result
}

type CommaAndParameter struct {
	*Comma
	*Parameter
}

func NewCommaAndParameter(parameter *Parameter) *CommaAndParameter {
	return &CommaAndParameter{
		Comma:     ConstComma,
		Parameter: parameter,
	}
}

func NewCommaAndParameterByToken(varTypeToken *Token, varNameToken *Token) *CommaAndParameter {
	return NewCommaAndParameter(NewParameterByToken(varTypeToken, varNameToken))
}

func (c *CommaAndParameter) ToXML() []string {
	result := []string{}
	result = append(result, c.Comma.ToXML())
	result = append(result, c.Parameter.ToXML()...)
	return result
}

type Parameter struct {
	*VarType
	*VarName
}

func NewParameter(varType *VarType, varName *VarName) *Parameter {
	return &Parameter{
		VarType: varType,
		VarName: varName,
	}
}

func NewParameterByToken(varTypeToken *Token, varNameToken *Token) *Parameter {
	return NewParameter(NewVarType(varTypeToken), NewVarName(varNameToken))
}

func (p *Parameter) Check() error {
	if err := p.VarType.Check(); err != nil {
		return err
	}

	if err := p.VarName.Check(); err != nil {
		return err
	}

	return nil
}

func (p *Parameter) ToXML() []string {
	result := []string{}
	result = append(result, p.VarType.ToXML())
	result = append(result, p.VarName.ToXML())
	return result
}

// '{' varDec* statements* '}'
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
		subroutineBody.VarDecs.Add(varDec)
	}

	// TODO Statements

	// TODO subroutineBodyの他の部分の実装が終わったら有効にする
	//closingCurlyBracket := p.advanceToken()
	//if err := ConstClosingCurlyBracket.Check(closingCurlyBracket); err != nil {
	//	return nil, err
	//}

	//fmt.Println(p.tokens.debug())
	//fmt.Println(p.readFirstToken().debug())

	return subroutineBody, nil
}

type SubroutineBody struct {
	*VarDecs
	*Statements
	*OpeningCurlyBracket
	*ClosingCurlyBracket
}

func NewSubroutineBody() *SubroutineBody {
	return &SubroutineBody{
		VarDecs:             NewVarDecs(),
		OpeningCurlyBracket: ConstOpeningCurlyBracket,
		ClosingCurlyBracket: ConstClosingCurlyBracket,
	}
}

func (s *SubroutineBody) ToXML() []string {
	result := []string{}
	result = append(result, "<subroutineBody>")
	result = append(result, s.OpeningCurlyBracket.ToXML())
	result = append(result, s.VarDecs.ToXML()...)
	result = append(result, s.ClosingCurlyBracket.ToXML())
	result = append(result, "</subroutineBody>")
	return result
}

type VarDecs struct {
	Items []*VarDec
}

func NewVarDecs() *VarDecs {
	return &VarDecs{
		Items: []*VarDec{},
	}
}

func (v *VarDecs) Add(item *VarDec) {
	v.Items = append(v.Items, item)
}

func (v *VarDecs) IsVarDecKeyword(token *Token) bool {
	if token == nil {
		return false
	}
	return token.Value == NewVarDec().Keyword.Value
}

func (v *VarDecs) ToXML() []string {
	result := []string{}
	for _, item := range v.Items {
		result = append(result, item.ToXML()...)
	}
	return result
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

type VarDec struct {
	*Keyword
	*VarType
	*VarNames
	*Semicolon
}

func NewVarDec() *VarDec {
	return &VarDec{
		Keyword:   NewKeywordByValue("var"),
		VarNames:  NewVarNames(),
		Semicolon: ConstSemicolon,
	}
}

func (v *VarDec) CheckKeyword(token *Token) error {
	return NewKeyword(token).Check(v.Keyword.Value)
}

func (v *VarDec) SetVarType(token *Token) error {
	varType := NewVarType(token)
	if err := varType.Check(); err != nil {
		return err
	}

	v.VarType = varType
	return nil
}

func (v *VarDec) ToXML() []string {
	result := []string{}
	result = append(result, "<varDec>")
	result = append(result, v.Keyword.ToXML())
	result = append(result, v.VarType.ToXML())
	result = append(result, v.VarNames.ToXML()...)
	result = append(result, v.Semicolon.ToXML())
	result = append(result, "</varDec>")
	return result
}

type Statements struct {
	*NotImplemented
}

type VarType struct {
	*Token
}

func NewVarType(token *Token) *VarType {
	return &VarType{
		Token: token,
	}
}

func (v *VarType) IsCheck() bool {
	err := v.Check()
	return err == nil
}

func (v *VarType) Check() error {
	if err := v.CheckIdentifier(); err == nil {
		return nil
	}

	expected := []string{"int", "char", "boolean"}
	if err := v.CheckKeywordValue(expected...); err == nil {
		return nil
	}

	message := fmt.Sprintf("VarType: got = %s", v.debug())
	return errors.New(message)
}

type VarNames struct {
	First            *VarName
	CommaAndVarNames []*CommaAndVarName
}

func NewVarNames() *VarNames {
	return &VarNames{
		CommaAndVarNames: []*CommaAndVarName{},
	}
}

func (v *VarNames) SetFirstVarName(token *Token) error {
	varName := NewVarName(token)
	if err := varName.Check(); err != nil {
		return err
	}

	v.First = varName
	return nil
}

func (v *VarNames) AddCommaAndVarName(commaToken *Token, varNameToken *Token) error {
	if err := ConstComma.Check(commaToken); err != nil {
		return err
	}

	varName := NewVarName(varNameToken)
	if err := varName.Check(); err != nil {
		return err
	}

	v.CommaAndVarNames = append(v.CommaAndVarNames, NewCommaAndVarName(varName))
	return nil
}

func (v *VarNames) ToXML() []string {
	result := []string{}
	result = append(result, v.First.ToXML())
	for _, commaAndVarName := range v.CommaAndVarNames {
		result = append(result, commaAndVarName.ToXML()...)
	}
	return result
}

type CommaAndVarName struct {
	Comma   *Comma
	VarName *VarName
}

func NewCommaAndVarName(varName *VarName) *CommaAndVarName {
	return &CommaAndVarName{
		Comma:   ConstComma,
		VarName: varName,
	}
}

func NewCommaAndVarNameByValue(value string) *CommaAndVarName {
	return NewCommaAndVarName(NewVarNameByValue(value))
}

func (c *CommaAndVarName) ToXML() []string {
	return []string{
		c.Comma.ToXML(),
		c.VarName.ToXML(),
	}
}

type SubroutineName struct {
	*Identifier
}

func NewSubroutineName(token *Token) *SubroutineName {
	return &SubroutineName{
		Identifier: NewIdentifier("SubroutineName", token),
	}
}

func NewSubroutineNameByValue(value string) *SubroutineName {
	return NewSubroutineName(NewToken(value, TokenIdentifier))
}

type VarName struct {
	*Identifier
}

func NewVarName(token *Token) *VarName {
	return &VarName{
		Identifier: NewIdentifier("VarName", token),
	}
}

func NewVarNameByValue(value string) *VarName {
	return NewVarName(NewToken(value, TokenIdentifier))
}

type Keyword struct {
	*Token
}

func NewKeyword(token *Token) *Keyword {
	return &Keyword{
		Token: token,
	}
}

func NewKeywordByValue(value string) *Keyword {
	return NewKeyword(NewToken(value, TokenKeyword))
}

func (k *Keyword) Check(expected ...string) error {
	return k.CheckKeywordValue(expected...)
}

type Symbol struct {
	*Token
}

func NewSymbol(token *Token) *Symbol {
	return &Symbol{
		Token: token,
	}
}

func NewSymbolByValue(value string) *Symbol {
	return NewSymbol(NewToken(value, TokenSymbol))
}

func (s *Symbol) Check(expected string) error {
	return s.CheckSymbolValue(expected)
}

type Identifier struct {
	Name string
	*Token
}

func NewIdentifier(name string, token *Token) *Identifier {
	return &Identifier{
		Name:  name,
		Token: token,
	}
}

func (i *Identifier) Check() error {
	return i.CheckTokenType(TokenIdentifier, "Identifier "+i.Name)
}

// よく使われるシンボル
// [] - Square brackets
var ConstOpeningCurlyBracket = NewOpeningCurlyBracket()
var ConstClosingCurlyBracket = NewClosingCurlyBracket()
var ConstOpeningRoundBracket = NewOpeningRoundBracket()
var ConstClosingRoundBracket = NewClosingRoundBracket()
var ConstComma = NewComma()
var ConstSemicolon = NewSemicolon()

type OpeningCurlyBracket struct {
	*Symbol
}

func NewOpeningCurlyBracket() *OpeningCurlyBracket {
	return &OpeningCurlyBracket{
		Symbol: NewSymbolByValue("{"),
	}
}

func (o *OpeningCurlyBracket) Check(token *Token) error {
	return NewSymbol(token).Check(o.Value)
}

type ClosingCurlyBracket struct {
	*Symbol
}

func NewClosingCurlyBracket() *ClosingCurlyBracket {
	return &ClosingCurlyBracket{
		Symbol: NewSymbolByValue("}"),
	}
}

func (c *ClosingCurlyBracket) Check(token *Token) error {
	return NewSymbol(token).Check(c.Value)
}

type OpeningRoundBracket struct {
	*Symbol
}

func NewOpeningRoundBracket() *OpeningRoundBracket {
	return &OpeningRoundBracket{
		Symbol: NewSymbolByValue("("),
	}
}

func (o *OpeningRoundBracket) Check(token *Token) error {
	return NewSymbol(token).Check(o.Value)
}

type ClosingRoundBracket struct {
	*Symbol
}

func NewClosingRoundBracket() *ClosingRoundBracket {
	return &ClosingRoundBracket{
		Symbol: NewSymbolByValue(")"),
	}
}

func (c *ClosingRoundBracket) Check(token *Token) error {
	return NewSymbol(token).Check(c.Value)
}

type Comma struct {
	*Symbol
}

func NewComma() *Comma {
	return &Comma{
		Symbol: NewSymbolByValue(","),
	}
}

func (c *Comma) IsCheck(token *Token) bool {
	err := c.Check(token)
	return err == nil
}

func (c *Comma) Check(token *Token) error {
	return NewSymbol(token).Check(c.Value)
}

type Semicolon struct {
	*Symbol
}

func NewSemicolon() *Semicolon {
	return &Semicolon{
		Symbol: NewSymbolByValue(";"),
	}
}

func (s *Semicolon) Check(token *Token) error {
	return NewSymbol(token).Check(s.Value)
}

type NotImplemented struct {
	*Token
}
