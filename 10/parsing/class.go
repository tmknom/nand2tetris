package parsing

import "../token"

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

func (c *Class) CheckKeyword(token *token.Token) error {
	return NewKeyword(token).Check(c.Keyword.Value)
}

type ClassName struct {
	*Identifier
}

func NewClassName(token *token.Token) *ClassName {
	return &ClassName{
		Identifier: NewIdentifier("ClassName", token),
	}
}

func (c *Class) SetClassName(token *token.Token) error {
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

func (c *ClassVarDecs) hasClassVarDec(token *token.Token) bool {
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

func (c *ClassVarDec) SetKeyword(token *token.Token) error {
	if err := c.checkKeyword(token); err != nil {
		return err
	}

	c.Keyword = NewKeywordByValue(token.Value)
	return nil
}

func (c *ClassVarDec) checkKeyword(token *token.Token) error {
	expected := []string{"static", "field"}
	return token.CheckKeywordValue(expected...)
}

func (c *ClassVarDec) SetVarType(token *token.Token) error {
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
