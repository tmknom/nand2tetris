package main

import "fmt"

type SymbolTable struct {
	entries     map[string]int
	nextAddress int
}

const InitialAddress = 16

func NewSymbolTable() *SymbolTable {
	entries := map[string]int{}
	entries["SP"] = 0
	entries["LCL"] = 1
	entries["ARG"] = 2
	entries["THIS"] = 3
	entries["THAT"] = 4
	entries["SCREEN"] = 16384
	entries["KBD"] = 24576

	for i := 0; i < 16; i++ {
		key := fmt.Sprintf("R%d", i)
		entries[key] = i
	}

	return &SymbolTable{entries: entries, nextAddress: InitialAddress}
}

func (s *SymbolTable) AddEntry(symbol string, address int) {
	s.entries[symbol] = address
}

func (s *SymbolTable) AddVariableEntry(symbol string) int {
	s.entries[symbol] = s.nextAddress
	s.nextAddress += 1
	return s.entries[symbol]
}

func (s *SymbolTable) Address(symbol string) int {
	result, ok := s.entries[symbol]
	if !ok {
		result = s.AddVariableEntry(symbol)
	}
	return result
}
