package symbol

import "fmt"

var GlobalIdGenerator = &IdGenerator{id: 0}

type IdGenerator struct {
	id int
}

func (i *IdGenerator) Generate() string {
	i.id += 1
	return fmt.Sprintf("ID_%d", i.id)
}

func (i *IdGenerator) Reset() {
	i.id = 0
}
