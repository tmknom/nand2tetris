package symbol

import "fmt"

const DebugSymbolTables = false

var GlobalSymbolTables = NewSymbolTables("Uninitialized")

type SymbolTables struct {
	*ClassSymbolTable
	*SubroutineSymbolTable
}

func NewSymbolTables(className string) *SymbolTables {
	return &SymbolTables{
		ClassSymbolTable:      NewClassSymbolTable(className),
		SubroutineSymbolTable: NewSubroutineSymbolTable("Uninitialized"),
	}
}

func (s *SymbolTables) Reset(className string) {
	s.ClassSymbolTable = NewClassSymbolTable(className)
	s.SubroutineSymbolTable = NewSubroutineSymbolTable("Uninitialized")
}

func (s *SymbolTables) ResetSubroutine(subroutineName string) {
	s.SubroutineSymbolTable = NewSubroutineSymbolTable(subroutineName)
}

func (s *SymbolTables) PrintClassSymbolTable() {
	if DebugSymbolTables {
		fmt.Println(s.ClassSymbolTable.String())
	}
}

func (s *SymbolTables) PrintSubroutineSymbolTable() {
	if DebugSymbolTables {
		fmt.Println(s.SubroutineSymbolTable.String())
	}
}

type SymbolTable struct {
	Items     []*SymbolItem
	Name      string
	TableType string
}

func NewSymbolTable(name string, tableType string) *SymbolTable {
	return &SymbolTable{
		Items:     []*SymbolItem{},
		Name:      name,
		TableType: tableType,
	}
}

func (s *SymbolTable) Add(item *SymbolItem) {
	s.Items = append(s.Items, item)
}

func (s *SymbolTable) String() string {
	result := fmt.Sprintf("%s = &%sSymbolTable{\n", s.Name, s.TableType)
	for i, item := range s.Items {
		result += fmt.Sprintf("  [%d] = %s\n", i, item.String())
	}
	result += "}\n"
	return result
}
