package symbol

type SubroutineSymbolTable struct {
	*SymbolTable
	*SubroutineScopeIndexer
}

func NewSubroutineSymbolTable(name string) *SubroutineSymbolTable {
	return &SubroutineSymbolTable{
		SymbolTable:            NewSymbolTable(name, "Subroutine"),
		SubroutineScopeIndexer: NewSubroutineScopeIndexer(),
	}
}

func (s *SubroutineSymbolTable) AddArgSymbol(name string, symbolType string) {
	scope := NewSymbolScope(ArgScope, s.SubroutineScopeIndexer.argIndex())
	item := NewSymbolItem(name, symbolType, scope, DefinedSymbol)
	s.Add(item)
}

func (s *SubroutineSymbolTable) AddVarSymbol(name string, symbolType string) {
	scope := NewSymbolScope(VarScope, s.SubroutineScopeIndexer.varIndex())
	item := NewSymbolItem(name, symbolType, scope, DefinedSymbol)
	s.Add(item)
}

type SubroutineScopeIndexer struct {
	ArgIndex int
	VarIndex int
}

func NewSubroutineScopeIndexer() *SubroutineScopeIndexer {
	return &SubroutineScopeIndexer{
		ArgIndex: 0,
		VarIndex: 0,
	}
}

func (s *SubroutineScopeIndexer) argIndex() int {
	result := s.ArgIndex
	s.ArgIndex += 1
	return result
}

func (s *SubroutineScopeIndexer) varIndex() int {
	result := s.VarIndex
	s.VarIndex += 1
	return result
}
