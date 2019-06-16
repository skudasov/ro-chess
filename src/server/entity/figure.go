package entity

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/name5566/leaf/log"
)

type figure struct {
	// Type of figure for client deserialiazation by interface
	Type string
	// Figure Name
	Name string
	// Visual mark for debug
	VisualMark string
	// Player who owns this figure
	Owner string
	// Is figure movable from pool
	Movable bool
	// Is figure active (making moves every turn)
	Active bool
	// Previous coordinates
	PrevX int
	PrevY int
	// This turn coordinates
	X int
	Y int
	// Is figure dead or not
	Alive bool
	// Initiative determines order of applying skills / making combat / moving of figure
	Initiative int
	HP         int
	MP         int
	AttackMin  int
	AttackMax  int
	Armor      int
}

// OptOwner option
func OptOwner(owner string) func(m *figure) {
	return func(m *figure) {
		m.Owner = owner
	}
}

// OptX option
func OptX(x int) func(m *figure) {
	return func(m *figure) {
		m.X = x
	}
}

// OptY option
func OptY(y int) func(m *figure) {
	return func(m *figure) {
		m.Y = y
	}
}

// OptPrevX option
func OptPrevX(x int) func(m *figure) {
	return func(m *figure) {
		m.PrevX = x
	}
}

// OptPrevY option
func OptPrevY(y int) func(m *figure) {
	return func(m *figure) {
		m.PrevY = y
	}
}

// OptAttackMin option
func OptAttackMin(a int) func(m *figure) {
	return func(m *figure) {
		m.AttackMin = a
	}
}

// OptAttackMax option
func OptAttackMax(a int) func(m *figure) {
	return func(m *figure) {
		m.AttackMax = a
	}
}

// OptDefence option
func OptDefence(d int) func(m *figure) {
	return func(m *figure) {
		m.Armor = d
	}
}

// OptMovable option
func OptMovable(movable bool) func(m *figure) {
	return func(m *figure) {
		m.Movable = movable
	}
}

// OptActive option
func OptActive(active bool) func(m *figure) {
	return func(m *figure) {
		m.Active = active
	}
}

// OptAlive option
func OptAlive(alive bool) func(m *figure) {
	return func(m *figure) {
		m.Alive = alive
	}
}

// OptHP option
func OptHP(hp int) func(m *figure) {
	return func(m *figure) {
		m.HP = hp
	}
}

// OptMP option
func OptMP(mp int) func(m *figure) {
	return func(m *figure) {
		m.MP = mp
	}
}

// PerformAttack calculates unit dmg
func (m *figure) PerformAttack() int {
	rand.Seed(time.Now().UnixNano())
	if m.AttackMax == m.AttackMin {
		return m.AttackMax
	}
	return rand.Intn(m.AttackMax-m.AttackMin) + m.AttackMin
}

// GetName gets unit name
func (m *figure) GetName() string {
	return m.Name
}

// GetVisualMark gets mark for visualization
func (m *figure) GetVisualMark() string {
	return m.VisualMark
}

// SetHP sets unit hp
func (m *figure) SetHP(hp int) {
	m.HP = hp
	if m.HP <= 0 {
		log.Debug("figure %s on %d, %d is dead", m.Name, m.X, m.Y)
		m.SetAlive(false)
	}
}

// GetHP gets unit hp
func (m *figure) GetHP() int {
	return m.HP
}

// SetMP sets unit mp
func (m *figure) SetMP(mp int) {
	m.MP = mp
}

// GetMP gets unit mp
func (m *figure) GetMP() int {
	return m.MP
}

// GetAttack gets unit min-max attack
func (m *figure) GetAttack() (int, int) {
	return m.AttackMin, m.AttackMax
}

// GetAttackStr gets unit min-max attack for visualization
func (m *figure) GetAttackStr() string {
	return fmt.Sprintf("%d-%d", m.AttackMin, m.AttackMax)
}

// GetDefence gets unit defences
func (m *figure) GetDefence() int {
	return m.Armor
}

// GetOwnerName gets unit owner
func (m *figure) GetOwnerName() string {
	return m.Owner
}

// GetMovable gets unit movable flag
func (m *figure) GetMovable() bool {
	return m.Movable
}

// SetMovable sets unit movable flag
func (m *figure) SetMovable(movable bool) {
	m.Movable = movable
}

// Activate activates unit making it walk forward for one cell
func (m *figure) Activate(walking bool) {
	m.Active = walking
}

// GetActive gets active flag
func (m *figure) GetActive() bool {
	return m.Active
}

// SetCoords sets unit X,Y coords
func (m *figure) SetCoords(X, Y int) {
	m.X = X
	m.Y = Y
}

// GetCoords gets unit X,Y coords
func (m *figure) GetCoords() (int, int) {
	return m.X, m.Y
}

// SetPrevCoords sets previous turn unit coords, needed for front-end to delete figure
func (m *figure) SetPrevCoords(X, Y int) {
	m.PrevX = X
	m.PrevY = Y
}

// GetPrevCoords gets previous turn unit coords
func (m *figure) GetPrevCoords() (int, int) {
	return m.PrevX, m.PrevY
}

// GetAlive gets unit alive flag
func (m *figure) GetAlive() bool {
	return m.Alive
}

// SetAlive sets unit alive flag
func (m *figure) SetAlive(alive bool) {
	m.Alive = alive
}

// GetInitiative gets unit initiative
func (m *figure) GetInitiative() int {
	return m.Initiative
}

// AddAttack adds min/max attack according to buff mechanics
func (m *figure) AddAttack(minAtk, maxAtk int) {
	m.AttackMin += minAtk
	m.AttackMax += maxAtk
}

// AddDefense adds defense according to buff mechanics
func (m *figure) AddDefense(def int) {
	m.Armor += def
}
