package symbol

import "fmt"

type SymbolItem struct {
	*SymbolName
	*SymbolType
	*SymbolScope
}

func NewSymbolItem(name string, symbolType string, symbolScope *SymbolScope) *SymbolItem {
	return &SymbolItem{
		SymbolName:  NewSymbolName(name),
		SymbolType:  NewSymbolType(symbolType),
		SymbolScope: symbolScope,
	}
}

func (s *SymbolItem) ToCode() string {
	return fmt.Sprintf("%s %d", s.ScopeKind, s.ScopeIndex)
}

func (s *SymbolItem) String() string {
	result := "&SymbolItem{ "
	result += fmt.Sprintf("Name: %s, ", s.SymbolName.Value)
	result += fmt.Sprintf("Type: %s, ", s.SymbolType.Value)
	result += fmt.Sprintf("Kind: %s, ", s.ScopeKind)
	result += fmt.Sprintf("Index: %d", s.ScopeIndex)
	result += " }"
	return result
}

type SymbolName struct {
	Value string
}

func NewSymbolName(value string) *SymbolName {
	return &SymbolName{
		Value: value,
	}
}

type SymbolType struct {
	Value string
}

func NewSymbolType(value string) *SymbolType {
	return &SymbolType{
		Value: value,
	}
}

type SymbolScope struct {
	ScopeKind
	ScopeIndex int
}

func NewSymbolScope(scopeKind ScopeKind, scopeIndex int) *SymbolScope {
	return &SymbolScope{
		ScopeKind:  scopeKind,
		ScopeIndex: scopeIndex,
	}
}

type ScopeKind int

const (
	_ ScopeKind = iota
	StaticScope
	FieldScope
	ArgScope
	VarScope
	ClassScope
	NoneScope
)

func (s ScopeKind) String() string {
	switch s {
	case StaticScope:
		return "static"
	case FieldScope:
		return "field"
	case ArgScope:
		return "argument"
	case VarScope:
		return "local"
	case ClassScope:
		return "class"
	case NoneScope:
		return "none"
	default:
		return "invalid ScopeKind"
	}
}
