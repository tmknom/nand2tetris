package parsing

import "../token"

type ParameterList struct {
	First              *Parameter
	CommaAndParameters []*CommaAndParameter
}

func NewParameterList() *ParameterList {
	return &ParameterList{
		CommaAndParameters: []*CommaAndParameter{},
	}
}

func (p *ParameterList) Add(varTypeToken *token.Token, varNameToken *token.Token) error {
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

func NewCommaAndParameterByToken(varTypeToken *token.Token, varNameToken *token.Token) *CommaAndParameter {
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

func NewParameterByToken(varTypeToken *token.Token, varNameToken *token.Token) *Parameter {
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
