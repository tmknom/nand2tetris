package parsing

import (
	"../token"
	"fmt"
	"github.com/pkg/errors"
)

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

func (s *SubroutineDecs) hasSubroutineDec(token *token.Token) bool {
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

func (s *SubroutineDec) SetSubroutineType(token *token.Token) error {
	subroutineType := NewSubroutineType(token)
	if err := subroutineType.Check(); err != nil {
		return err
	}

	s.SubroutineType = subroutineType
	return nil
}

func (s *SubroutineDec) SetSubroutineName(token *token.Token) error {
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
	*token.Token
}

func NewSubroutineType(token *token.Token) *SubroutineType {
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

	message := fmt.Sprintf("SubroutineType: got = %s", s.Debug())
	return errors.New(message)
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
		Statements:          NewStatements(),
		OpeningCurlyBracket: ConstOpeningCurlyBracket,
		ClosingCurlyBracket: ConstClosingCurlyBracket,
	}
}

func (s *SubroutineBody) ToXML() []string {
	result := []string{}
	result = append(result, "<subroutineBody>")
	result = append(result, s.OpeningCurlyBracket.ToXML())
	result = append(result, s.VarDecs.ToXML()...)
	result = append(result, s.Statements.ToXML()...)
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

func (v *VarDecs) AddVarDec(item *VarDec) {
	v.Items = append(v.Items, item)
}

func (v *VarDecs) IsVarDecKeyword(token *token.Token) bool {
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

func (v *VarDec) CheckKeyword(token *token.Token) error {
	return NewKeyword(token).Check(v.Keyword.Value)
}

func (v *VarDec) SetVarType(token *token.Token) error {
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
