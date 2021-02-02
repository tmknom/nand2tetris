package symbol

type ScopeIndexer struct {
	StaticIndex int
	FieldIndex  int
}

func NewIndexes() *ScopeIndexer {
	return &ScopeIndexer{
		StaticIndex: 0,
		FieldIndex:  0,
	}
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
