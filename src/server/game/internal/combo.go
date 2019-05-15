package internal

//go:generate stringer -type=comboType

import "fmt"

type comboType int

const (
	attack3Combo comboType = iota
	defense3Combo
)

type comboVal struct {
	UnitCoord []int
	Stacks    int
}
type comboStore map[comboKey]*comboVal

type comboKey struct {
	Axis   string
	Class  string
	Indent int
}

func (m comboKey) String() string {
	return fmt.Sprintf("%s-%s-%d", m.Axis, m.Class, m.Indent)
}
