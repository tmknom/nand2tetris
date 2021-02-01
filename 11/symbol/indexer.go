package symbol

type ScopeIndexer struct {
	StaticIndex int
}

func NewIndexes() *ScopeIndexer {
	return &ScopeIndexer{
		StaticIndex: 0,
	}
}

func (s *ScopeIndexer) staticIndex() int {
	result := s.StaticIndex
	s.StaticIndex += 1
	return result
}
