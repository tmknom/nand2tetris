package symbol

import (
	"fmt"
)

type SymbolTable struct {
	Name  string
	Items []*SymbolItem
	*ScopeIndexer
}

func NewSymbolTable(name string) *SymbolTable {
	return &SymbolTable{
		Name:         name,
		ScopeIndexer: NewIndexes(),
	}
}

func (s *SymbolTable) AddDefinedStaticSymbol(name string, symbolType string) {
	scope := NewSymbolScope(StaticScope, s.ScopeIndexer.staticIndex())
	item := NewSymbolItem(name, symbolType, scope, DefinedSymbol)
	s.Add(item)
}

func (s *SymbolTable) AddDefinedClassSymbol(name string) {
	scope := NewSymbolScope(ClassScope, 0)
	item := NewSymbolItem(name, name, scope, DefinedSymbol)
	s.Add(item)
}

func (s *SymbolTable) Add(item *SymbolItem) {
	s.Items = append(s.Items, item)
}

func (s *SymbolTable) String() string {
	result := "\n"
	result += fmt.Sprintf("%s = &SymbolTable{\n", s.Name)
	for i, item := range s.Items {
		result += fmt.Sprintf("  [%d] = %s\n", i, item.String())
	}
	result += "}\n"
	return result
}
