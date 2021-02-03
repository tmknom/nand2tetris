package symbol

import (
	"fmt"
	"github.com/pkg/errors"
)

var DebugSymbolTables = true

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

func (s *SymbolTables) Find(name string) (string, error) {
	subroutineSymbolItem, err := s.SubroutineSymbolTable.Find(name)
	if err == nil {
		return subroutineSymbolItem.ToCode(), nil
	}

	classSymbolItem, err := s.ClassSymbolTable.Find(name)
	if err == nil {
		return classSymbolItem.ToCode(), nil
	}

	message := fmt.Sprintf("not found at %s and %s: name = %s", s.SubroutineSymbolTable.TableName(), s.ClassSymbolTable.TableName(), name)
	return "", errors.New(message)
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

func (s *SymbolTable) Find(name string) (*SymbolItem, error) {
	for _, item := range s.Items {
		if item.SymbolName.Value == name {
			return item, nil
		}
	}

	message := fmt.Sprintf("not found at %s: name = %s", s.TableName(), name)
	return nil, errors.New(message)
}

func (s *SymbolTable) TableName() string {
	return fmt.Sprintf("%sSymbolTable(%s)", s.TableType, s.Name)
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
