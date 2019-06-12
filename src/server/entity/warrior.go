package entity

import (
	"fmt"
	"github.com/name5566/leaf/log"
	"math/rand"
	"time"
)

// Warrior unit
type Warrior struct {
	Figure   *figure
	SkillSet *SkillSet
}

// NewWarrior creates new peon unit
func NewWarrior(opts ...func(m *figure)) *Warrior {
	p := &Warrior{Figure: &figure{}}
	p.Figure.Type = "warrior"
	p.Figure.Name = "Warrior"
	p.Figure.VisualMark = "W"
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

// NewConstDmgWarrior creates new peon unit with const dps
func NewConstDmgWarrior(opts ...func(m *figure)) *Warrior {
	p := &Warrior{Figure: &figure{}}
	p.Figure.Type = "warrior"
	p.Figure.Name = "Warrior"
	p.Figure.VisualMark = "W"
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

// GetSkillSet Gets skillset
func (m *Warrior) GetSkillSet() *SkillSet {
	return m.SkillSet
}

// SetSkillSet Sets skill set
func (m *Warrior) SetSkillSet(ss *SkillSet) {
	m.SkillSet = ss
}

// GetRotation Gets skill rotation
func (m *Warrior) GetRotation() []*AppliedSkill {
	return m.SkillSet.Rotation
}

// LearnSkill Learns skill consuming xp
func (m *Warrior) LearnSkill(name string, skill SkillFunc) {
	if m.SkillSet == nil {
		m.SkillSet = NewEmptySkillSet()
	}
	m.SkillSet.SkillBook[name] = skill
}

func (m *Warrior) AddSkillToRotation(boardName string, skillName string, from Pair, to Pair) {
	log.Debug("adding skill to rotation: %s: %d, %d -> %d, %d", skillName, from.X, from.Y, to.X, to.Y)
	m.SkillSet.Rotation = append(m.SkillSet.Rotation, &AppliedSkill{
		boardName,
		1,
		from,
		to,
		skillName,
	},
	)
}

// ApplySkills Applies skills from rotation in order, calling every skill from self skillbook
func (m *Warrior) ApplySkills(updatedFigures *[]Figurable, updatedPlayers *[]Player, clog *[]CombatEvent) {
	for _, app := range m.SkillSet.Rotation {
		if _, ok := m.SkillSet.SkillBook[app.Name]; !ok {
			log.Debug("no such skill in the book: %s", app.Name)
			continue
		}
		m.SkillSet.SkillBook[app.Name](app.Board, app.From, app.To, updatedFigures, updatedPlayers, clog)
	}
}

// PerformAttack calculates unit dmg
func (m *Warrior) PerformAttack() int {
	rand.Seed(time.Now().UnixNano())
	if m.Figure.AttackMax == m.Figure.AttackMin {
		return m.Figure.AttackMax
	}
	return rand.Intn(m.Figure.AttackMax-m.Figure.AttackMin) + m.Figure.AttackMin
}

// GetName gets unit name
func (m *Warrior) GetName() string {
	return m.Figure.Name
}

// GetVisualMark gets mark for visualization
func (m *Warrior) GetVisualMark() string {
	return m.Figure.VisualMark
}

// SetHP sets unit hp
func (m *Warrior) SetHP(hp int) {
	m.Figure.HP = hp
}

// GetHP gets unit hp
func (m *Warrior) GetHP() int {
	return m.Figure.HP
}

// SetMP sets unit mp
func (m *Warrior) SetMP(mp int) {
	m.Figure.MP = mp
}

// GetMP gets unit mp
func (m *Warrior) GetMP() int {
	return m.Figure.MP
}

// GetAttack gets unit min-max attack
func (m *Warrior) GetAttack() (int, int) {
	return m.Figure.AttackMin, m.Figure.AttackMax
}

// GetAttackStr gets unit min-max attack for visualization
func (m *Warrior) GetAttackStr() string {
	return fmt.Sprintf("%d-%d", m.Figure.AttackMin, m.Figure.AttackMax)
}

// GetDefence gets unit defences
func (m *Warrior) GetDefence() int {
	return m.Figure.Armor
}

// GetOwnerName gets unit owner
func (m *Warrior) GetOwnerName() string {
	return m.Figure.Owner
}

// GetMovable gets unit movable flag
func (m *Warrior) GetMovable() bool {
	return m.Figure.Movable
}

// SetMovable sets unit movable flag
func (m *Warrior) SetMovable(movable bool) {
	m.Figure.Movable = movable
}

// Activate activates unit making it walk forward for one cell
func (m *Warrior) Activate(walking bool) {
	m.Figure.Active = walking
}

// GetActive gets active flag
func (m *Warrior) GetActive() bool {
	return m.Figure.Active
}

// SetCoords sets unit X,Y coords
func (m *Warrior) SetCoords(X, Y int) {
	m.Figure.X = X
	m.Figure.Y = Y
}

// GetCoords gets unit X,Y coords
func (m *Warrior) GetCoords() (int, int) {
	return m.Figure.X, m.Figure.Y
}

// SetPrevCoords sets previous turn unit coords, needed for front-end to delete figure
func (m *Warrior) SetPrevCoords(X, Y int) {
	m.Figure.PrevX = X
	m.Figure.PrevY = Y
}

// GetPrevCoords gets previous turn unit coords
func (m *Warrior) GetPrevCoords() (int, int) {
	return m.Figure.PrevX, m.Figure.PrevY
}

// GetAlive gets unit alive flag
func (m *Warrior) GetAlive() bool {
	return m.Figure.Alive
}

// SetAlive sets unit alive flag
func (m *Warrior) SetAlive(alive bool) {
	m.Figure.Alive = alive
}

// Clone copies unit
func (m *Warrior) Clone() Figurable {
	cloned := *m
	return &cloned
}

// GetInitiative gets unit initiative
func (m *Warrior) GetInitiative() int {
	return m.Figure.Initiative
}

// AddAttack adds min/max attack according to buff mechanics
func (m *Warrior) AddAttack(minAtk, maxAtk int) {
	m.Figure.AttackMin += minAtk
	m.Figure.AttackMax += maxAtk
}

// AddDefense adds defense accorting to buff mechanics
func (m *Warrior) AddDefense(def int) {
	m.Figure.Armor += def
}
