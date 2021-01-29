package parsing

import "../token"

type Statements struct {
	Items []*Statement
}

func NewStatements() *Statements {
	return &Statements{
		Items: []*Statement{},
	}
}

func (s *Statements) AddStatement(item *Statement) {
	s.Items = append(s.Items, item)
}

func (s *Statements) ToXML() []string {
	result := []string{}
	result = append(result, "<statements>")
	for _, item := range s.Items {
		result = append(result, item.ToXML()...)
	}
	result = append(result, "</statements>")
	return result
}

func (s *Statements) IsStatementKeyword(token *token.Token) bool {
	if token == nil {
		return false
	}

	keywords := []string{
		"let",
		"if",
		"while",
		"do",
		"return",
	}

	err := NewKeyword(token).CheckKeywordValue(keywords...)
	return err == nil
}

type Statement struct {
	*NotImplemented
}

func NewStatement() *Statement {
	return &Statement{}
}

func (s *Statement) ToXML() []string {
	return []string{}
}
