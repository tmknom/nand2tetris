package symbol

import "fmt"

type SymbolItem struct {
	*SymbolName
	*SymbolType
	*SymbolScope
	SymbolStatus
}

func NewSymbolItem(name string, symbolType string, symbolScope *SymbolScope, symbolStatus SymbolStatus) *SymbolItem {
	return &SymbolItem{
		SymbolName:   NewSymbolName(name),
		SymbolType:   NewSymbolType(symbolType),
		SymbolScope:  symbolScope,
		SymbolStatus: symbolStatus,
	}
}

func (s *SymbolItem) String() string {
	result := "&SymbolItem{ "
	result += fmt.Sprintf("Name: %s, ", s.SymbolName.Value)
	result += fmt.Sprintf("Type: %s, ", s.SymbolType.Value)
	result += fmt.Sprintf("Scope: %s[%d], ", s.ScopeKind, s.ScopeIndex)
	result += fmt.Sprintf("Status: %s", s.SymbolStatus)
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
		return "arg"
	case VarScope:
		return "var"
	case ClassScope:
		return "class"
	case NoneScope:
		return "none"
	default:
		return "invalid ScopeKind"
	}
}

type SymbolStatus int

const (
	_ SymbolStatus = iota
	DefinedSymbol
	UsedSymbol
)

func (s SymbolStatus) String() string {
	switch s {
	case DefinedSymbol:
		return "defined"
	case UsedSymbol:
		return "used"
	default:
		return "invalid SymbolStatus"
	}
}
