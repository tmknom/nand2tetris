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
	item := NewSymbolItem(name, symbolType, scope)
	s.Add(item)
}

func (s *ClassSymbolTable) AddFieldSymbol(name string, symbolType string) {
	scope := NewSymbolScope(FieldScope, s.ClassScopeIndexer.fieldIndex())
	item := NewSymbolItem(name, symbolType, scope)
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

func (c *ClassScopeIndexer) fieldIndex() int {
	result := c.FieldIndex
	c.FieldIndex += 1
	return result
}

func (c *ClassScopeIndexer) staticIndex() int {
	result := c.StaticIndex
	c.StaticIndex += 1
	return result
}

func (c *ClassScopeIndexer) FieldLength() int {
	return c.FieldIndex
}

func (c *ClassScopeIndexer) StaticLength() int {
	return c.StaticIndex
}
