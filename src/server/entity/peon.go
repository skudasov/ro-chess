package entity

import (
	"fmt"
	"math/rand"
	"time"
)

// Peon unit
type Peon struct {
	Figure *figure
}

// NewPeon creates new peon unit
func NewPeon(opts ...func(m *figure)) *Peon {
	p := &Peon{Figure: &figure{}}
	p.Figure.Type = "peon"
	p.Figure.Name = "Peon"
	p.Figure.VisualMark = "P"
	p.Figure.Movable = true
	p.Figure.Active = false
	p.Figure.PrevX = 0
	p.Figure.PrevY = 0
	p.Figure.Alive = true
	p.Figure.HP = 10
	p.Figure.MP = 0
	p.Figure.AttackMin = 1
	p.Figure.AttackMax = 3
	p.Figure.Armor = 0
	p.Figure.Initiative = 1
	for _, o := range opts {
		o(p.Figure)
	}
	return p
}

// NewConstDmgPeon creates new peon unit with const dps
func NewConstDmgPeon(opts ...func(m *figure)) *Peon {
	p := &Peon{Figure: &figure{}}
	p.Figure.Type = "peon"
	p.Figure.Name = "Peon"
	p.Figure.VisualMark = "P"
	p.Figure.Movable = true
	p.Figure.Active = false
	p.Figure.PrevX = 0
	p.Figure.PrevY = 0
	p.Figure.Alive = true
	p.Figure.HP = 10
	p.Figure.MP = 0
	p.Figure.AttackMin = 5
	p.Figure.AttackMax = 5
	p.Figure.Armor = 0
	p.Figure.Initiative = 1
	for _, o := range opts {
		o(p.Figure)
	}
	return p
}

// PerformAttack calculates unit dmg
func (m *Peon) PerformAttack() int {
	rand.Seed(time.Now().UnixNano())
	if m.Figure.AttackMax == m.Figure.AttackMin {
		return m.Figure.AttackMax
	}
	return rand.Intn(m.Figure.AttackMax-m.Figure.AttackMin) + m.Figure.AttackMin
}

// GetName gets unit name
func (m *Peon) GetName() string {
	return m.Figure.Name
}

// GetVisualMark gets mark for visualization
func (m *Peon) GetVisualMark() string {
	return m.Figure.VisualMark
}

// SetHP sets unit hp
func (m *Peon) SetHP(hp int) {
	m.Figure.HP = hp
}

// GetHP gets unit hp
func (m *Peon) GetHP() int {
	return m.Figure.HP
}

// SetMP sets unit mp
func (m *Peon) SetMP(mp int) {
	m.Figure.MP = mp
}

// GetMP gets unit mp
func (m *Peon) GetMP() int {
	return m.Figure.MP
}

// GetAttack gets unit min-max attack
func (m *Peon) GetAttack() (int, int) {
	return m.Figure.AttackMin, m.Figure.AttackMax
}

// GetAttackStr gets unit min-max attack for visualization
func (m *Peon) GetAttackStr() string {
	return fmt.Sprintf("%d-%d", m.Figure.AttackMin, m.Figure.AttackMax)
}

// GetDefence gets unit defences
func (m *Peon) GetDefence() int {
	return m.Figure.Armor
}

// GetOwnerName gets unit owner
func (m *Peon) GetOwnerName() string {
	return m.Figure.Owner
}

// GetMovable gets unit movable flag
func (m *Peon) GetMovable() bool {
	return m.Figure.Movable
}

// SetMovable sets unit movable flag
func (m *Peon) SetMovable(movable bool) {
	m.Figure.Movable = movable
}

// Activate activates unit making it walk forward for one cell
func (m *Peon) Activate(walking bool) {
	m.Figure.Active = walking
}

// GetActive gets active flag
func (m *Peon) GetActive() bool {
	return m.Figure.Active
}

// SetCoords sets unit x,y coords
func (m *Peon) SetCoords(X, Y int) {
	m.Figure.X = X
	m.Figure.Y = Y
}

// GetCoords gets unit x,y coords
func (m *Peon) GetCoords() (int, int) {
	return m.Figure.X, m.Figure.Y
}

// SetPrevCoords sets previous turn unit coords, needed for front-end to delete figure
func (m *Peon) SetPrevCoords(X, Y int) {
	m.Figure.PrevX = X
	m.Figure.PrevY = Y
}

// GetPrevCoords gets previous turn unit coords
func (m *Peon) GetPrevCoords() (int, int) {
	return m.Figure.PrevX, m.Figure.PrevY
}

// GetAlive gets unit alive flag
func (m *Peon) GetAlive() bool {
	return m.Figure.Alive
}

// SetAlive sets unit alive flag
func (m *Peon) SetAlive(alive bool) {
	m.Figure.Alive = alive
}

// Clone copies unit
func (m *Peon) Clone() Figurable {
	cloned := *m
	return &cloned
}

// GetInitiative gets unit initiative
func (m *Peon) GetInitiative() int {
	return m.Figure.Initiative
}

// AddAttack adds min/max attack according to buff mechanics
func (m *Peon) AddAttack(minAtk, maxAtk int) {
	m.Figure.AttackMin += minAtk
	m.Figure.AttackMax += maxAtk
}

// AddDefense adds defense accorting to buff mechanics
func (m *Peon) AddDefense(def int) {
	m.Figure.Armor += def
}
