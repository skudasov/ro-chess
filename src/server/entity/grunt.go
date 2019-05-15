package entity

import (
	"fmt"
	"math/rand"
	"time"
)

// Grunt unit
type Grunt struct {
	Figure *figure
}

// NewGrunt creates new grunt unit
func NewGrunt(opts ...func(m *figure)) *Grunt {
	g := &Grunt{Figure: &figure{}}
	g.Figure.Type = "grunt"
	g.Figure.Name = "Grunt"
	g.Figure.VisualMark = "G"
	g.Figure.Movable = true
	g.Figure.Active = false
	g.Figure.PrevX = 0
	g.Figure.PrevY = 0
	g.Figure.Alive = true
	g.Figure.HP = 30
	g.Figure.MP = 0
	g.Figure.AttackMin = 7
	g.Figure.AttackMax = 10
	g.Figure.Armor = 1
	for _, o := range opts {
		o(g.Figure)
	}
	return g
}

// PerformAttack calculates unit dmg
func (m *Grunt) PerformAttack() int {
	rand.Seed(time.Now().UnixNano())
	if m.Figure.AttackMax == m.Figure.AttackMin {
		return m.Figure.AttackMax
	}
	return rand.Intn(m.Figure.AttackMax-m.Figure.AttackMin) + m.Figure.AttackMin
}

// GetName gets unit name
func (m *Grunt) GetName() string {
	return m.Figure.Name
}

// GetVisualMark gets mark for visualization
func (m *Grunt) GetVisualMark() string {
	return m.Figure.VisualMark
}

// SetHP sets unit hp
func (m *Grunt) SetHP(hp int) {
	m.Figure.HP = hp
}

// GetHP gets unit hp
func (m *Grunt) GetHP() int {
	return m.Figure.HP
}

// SetMP sets unit mp
func (m *Grunt) SetMP(mp int) {
	m.Figure.MP = mp
}

// GetMP gets unit mp
func (m *Grunt) GetMP() int {
	return m.Figure.MP
}

// GetAttack gets unit min-max attack
func (m *Grunt) GetAttack() (int, int) {
	return m.Figure.AttackMin, m.Figure.AttackMax
}

// GetAttackStr gets unit min-max attack for visualization
func (m *Grunt) GetAttackStr() string {
	return fmt.Sprintf("%d-%d", m.Figure.AttackMin, m.Figure.AttackMax)
}

// GetDefence gets unit defences
func (m *Grunt) GetDefence() int {
	return m.Figure.Armor
}

// GetOwner gets unit owner
func (m *Grunt) GetOwner() string {
	return m.Figure.Owner
}

// GetMovable gets unit movable flag
func (m *Grunt) GetMovable() bool {
	return m.Figure.Movable
}

// SetMovable sets unit movable flag
func (m *Grunt) SetMovable(movable bool) {
	m.Figure.Movable = movable
}

// Activate activates unit making it walk forward for one cell
func (m *Grunt) Activate(walking bool) {
	m.Figure.Active = walking
}

// GetActive gets active flag
func (m *Grunt) GetActive() bool {
	return m.Figure.Active
}

// SetCoords sets unit x,y coords
func (m *Grunt) SetCoords(X, Y int) {
	m.Figure.X = X
	m.Figure.Y = Y
}

// GetCoords gets unit x,y coords
func (m *Grunt) GetCoords() (int, int) {
	return m.Figure.X, m.Figure.Y
}

// SetPrevCoords sets previous turn unit coords, needed for front-end to delete figure
func (m *Grunt) SetPrevCoords(X, Y int) {
	m.Figure.PrevX = X
	m.Figure.PrevY = Y
}

// GetPrevCoords gets previous turn unit coords
func (m *Grunt) GetPrevCoords() (int, int) {
	return m.Figure.PrevX, m.Figure.PrevY
}

// GetAlive gets unit alive flag
func (m *Grunt) GetAlive() bool {
	return m.Figure.Alive
}

// SetAlive sets unit alive flag
func (m *Grunt) SetAlive(alive bool) {
	m.Figure.Alive = alive
}

// Clone copies unit
func (m *Grunt) Clone() Figurable {
	figure := *m.Figure
	cloned := *m
	cloned.Figure = &figure
	return &cloned
}

// AddAttack adds min/max attack according to buff mechanics
func (m *Grunt) AddAttack(minAtk, maxAtk int) {
	m.Figure.AttackMin += minAtk
	m.Figure.AttackMax += maxAtk
}

// AddDefense adds defense accorting to buff mechanics
func (m *Grunt) AddDefense(def int) {
	m.Figure.Armor += def
}
