package internal

//go:generate stringer -type=side
//go:generate stringer -type=race

import (
	"github.com/f4hrenh9it/ro-chess/src/server/entity"
	"github.com/name5566/leaf/gate"
)

type side int

const (
	top side = iota
	bottom
)

func (m *side) Other() side {
	return bottom
}

type race int

const (
	orcs race = iota
	humans
)

type player struct {
	Name         string
	Agent        gate.Agent
	Opponent     *player
	Info         string
	Side         side
	Board        *board
	HP           int
	MP           int
	ZoneScoreY   int
	ComboStore   comboStore
	BoardFigures []entity.Figurable
	FigurePool   entity.Poolable
}

type figuresByInitiative []entity.Figurable

func (a figuresByInitiative) Len() int { return len(a) }
func (a figuresByInitiative) Less(i, j int) bool {
	return a[i].GetInitiative() > a[j].GetInitiative()
}
func (a figuresByInitiative) Swap(i, j int) { a[i], a[j] = a[j], a[i] }

type cell struct {
	Figure  entity.Figurable
	Terrain int
	Busy    bool
}

type board struct {
	Winner           string
	Turn             string
	TurnEnds         map[string]bool
	ZoneStartTopX    []int
	ZoneStartTopY    []int
	ZoneStartBottomX []int
	ZoneStartBottomY []int
	Players          map[string]*player
	Figures          []entity.Figurable
	Canvas           [][]*cell
}
