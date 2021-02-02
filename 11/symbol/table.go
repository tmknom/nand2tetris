package symbol

import "fmt"

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
