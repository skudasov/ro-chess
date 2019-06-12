package entity

import (
	"fmt"
	"github.com/name5566/leaf/log"
	"math/rand"
	"time"
)

// Mage unit
type Mage struct {
	Figure   *figure
	SkillSet *SkillSet
}

// NewMage creates new grunt unit
func NewMage(opts ...func(m *figure)) *Mage {
	g := &Mage{Figure: &figure{}}
	g.Figure.Type = "mage"
	g.Figure.Name = "Mage"
	g.Figure.VisualMark = "M"
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

// GetRotation Gets skills rotation
func (m *Mage) GetRotation() []*AppliedSkill {
	return m.SkillSet.Rotation
}

// AddSkillToRotation Adds skill to rotation
func (m *Mage) AddSkillToRotation(boardName string, skillName string, from Pair, to Pair) {
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

// GetSkillSet Gets skillset
func (m *Mage) GetSkillSet() *SkillSet {
	return m.SkillSet
}

// SetSkillSet Sets skill set
func (m *Mage) SetSkillSet(ss *SkillSet) {
	m.SkillSet = ss
}

// LearnSkill Learns skill consuming xp
func (m *Mage) LearnSkill(name string, skill SkillFunc) {
	//TODO: consume xp here?
	if m.SkillSet == nil {
		m.SkillSet = NewEmptySkillSet()
	}
	m.SkillSet.SkillBook[name] = skill
}

// ApplySkills Applies skills from rotation in order, calling every skill from self skillbook
func (m *Mage) ApplySkills(updatedFigures *[]Figurable, updatedPlayers *[]Player, clog *[]CombatEvent) {
	for _, app := range m.SkillSet.Rotation {
		if _, ok := m.SkillSet.SkillBook[app.Name]; !ok {
			log.Debug("no such skill in the book: %s", app.Name)
			continue
		}
		m.SkillSet.SkillBook[app.Name](app.Board, app.From, app.To, updatedFigures, updatedPlayers, clog)
	}
}

// PerformAttack calculates unit dmg
func (m *Mage) PerformAttack() int {
	rand.Seed(time.Now().UnixNano())
	if m.Figure.AttackMax == m.Figure.AttackMin {
		return m.Figure.AttackMax
	}
	return rand.Intn(m.Figure.AttackMax-m.Figure.AttackMin) + m.Figure.AttackMin
}

// GetName gets unit name
func (m *Mage) GetName() string {
	return m.Figure.Name
}

// GetVisualMark gets mark for visualization
func (m *Mage) GetVisualMark() string {
	return m.Figure.VisualMark
}

// SetHP sets unit hp
func (m *Mage) SetHP(hp int) {
	m.Figure.HP = hp
}

// GetHP gets unit hp
func (m *Mage) GetHP() int {
	return m.Figure.HP
}

// SetMP sets unit mp
func (m *Mage) SetMP(mp int) {
	m.Figure.MP = mp
}

// GetMP gets unit mp
func (m *Mage) GetMP() int {
	return m.Figure.MP
}

// GetAttack gets unit min-max attack
func (m *Mage) GetAttack() (int, int) {
	return m.Figure.AttackMin, m.Figure.AttackMax
}

// GetAttackStr gets unit min-max attack for visualization
func (m *Mage) GetAttackStr() string {
	return fmt.Sprintf("%d-%d", m.Figure.AttackMin, m.Figure.AttackMax)
}

// GetDefence gets unit defences
func (m *Mage) GetDefence() int {
	return m.Figure.Armor
}

// GetOwnerName gets unit owner
func (m *Mage) GetOwnerName() string {
	return m.Figure.Owner
}

// GetMovable gets unit movable flag
func (m *Mage) GetMovable() bool {
	return m.Figure.Movable
}

// SetMovable sets unit movable flag
func (m *Mage) SetMovable(movable bool) {
	m.Figure.Movable = movable
}

// Activate activates unit making it walk forward for one cell
func (m *Mage) Activate(walking bool) {
	m.Figure.Active = walking
}

// GetActive gets active flag
func (m *Mage) GetActive() bool {
	return m.Figure.Active
}

// SetCoords sets unit X,Y coords
func (m *Mage) SetCoords(X, Y int) {
	m.Figure.X = X
	m.Figure.Y = Y
}

// GetCoords gets unit X,Y coords
func (m *Mage) GetCoords() (int, int) {
	return m.Figure.X, m.Figure.Y
}

// SetPrevCoords sets previous turn unit coords, needed for front-end to delete figure
func (m *Mage) SetPrevCoords(X, Y int) {
	m.Figure.PrevX = X
	m.Figure.PrevY = Y
}

// GetPrevCoords gets previous turn unit coords
func (m *Mage) GetPrevCoords() (int, int) {
	return m.Figure.PrevX, m.Figure.PrevY
}

// GetAlive gets unit alive flag
func (m *Mage) GetAlive() bool {
	return m.Figure.Alive
}

// SetAlive sets unit alive flag
func (m *Mage) SetAlive(alive bool) {
	m.Figure.Alive = alive
}

// Clone copies unit
func (m *Mage) Clone() Figurable {
	figure := *m.Figure
	cloned := *m
	cloned.Figure = &figure
	return &cloned
}

// GetInitiative gets unit initiative
func (m *Mage) GetInitiative() int {
	return m.Figure.Initiative
}

// AddAttack adds min/max attack according to buff mechanics
func (m *Mage) AddAttack(minAtk, maxAtk int) {
	m.Figure.AttackMin += minAtk
	m.Figure.AttackMax += maxAtk
}

// AddDefense adds defense accorting to buff mechanics
func (m *Mage) AddDefense(def int) {
	m.Figure.Armor += def
}
