package symbol

type ScopeIndexer struct {
	StaticIndex int
	FieldIndex  int
	ArgIndex    int
	VarIndex    int
}

func NewIndexes() *ScopeIndexer {
	return &ScopeIndexer{
		StaticIndex: 0,
		FieldIndex:  0,
		ArgIndex:    0,
		VarIndex:    0,
	}
}

func (s *ScopeIndexer) varIndex() int {
	result := s.VarIndex
	s.VarIndex += 1
	return result
}

func (s *ScopeIndexer) argIndex() int {
	result := s.ArgIndex
	s.ArgIndex += 1
	return result
}

func (s *ScopeIndexer) fieldIndex() int {
	result := s.FieldIndex
	s.FieldIndex += 1
	return result
}

func (s *ScopeIndexer) staticIndex() int {
	result := s.StaticIndex
	s.StaticIndex += 1
	return result
}
