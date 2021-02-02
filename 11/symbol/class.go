package symbol

type ClassSymbolTable struct {
	*SymbolTable
	*ClassScopeIndexer
}

func NewClassSymbolTable(name string) *ClassSymbolTable {
	return &ClassSymbolTable{
		SymbolTable:       NewSymbolTable(name, "Class"),
		ClassScopeIndexer: NewClassScopeIndexer(),
	}
}

func (s *ClassSymbolTable) AddStaticSymbol(name string, symbolType string) {
	scope := NewSymbolScope(StaticScope, s.ClassScopeIndexer.staticIndex())
	item := NewSymbolItem(name, symbolType, scope, DefinedSymbol)
	s.Add(item)
}

func (s *ClassSymbolTable) AddFieldSymbol(name string, symbolType string) {
	scope := NewSymbolScope(FieldScope, s.ClassScopeIndexer.fieldIndex())
	item := NewSymbolItem(name, symbolType, scope, DefinedSymbol)
	s.Add(item)
}

type ClassScopeIndexer struct {
	StaticIndex int
	FieldIndex  int
}

func NewClassScopeIndexer() *ClassScopeIndexer {
	return &ClassScopeIndexer{
		StaticIndex: 0,
		FieldIndex:  0,
	}
}

func (s *ClassScopeIndexer) fieldIndex() int {
	result := s.FieldIndex
	s.FieldIndex += 1
	return result
}

func (s *ClassScopeIndexer) staticIndex() int {
	result := s.StaticIndex
	s.StaticIndex += 1
	return result
}
