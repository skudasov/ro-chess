package entity

import (
	"fmt"
	"github.com/f4hrenh9it/ro-chess/src/server/skills"
	"github.com/name5566/leaf/log"
	"math/rand"
	"time"
)

// Mage unit
type Mage struct {
	Figure   *figure
	SkillSet *skills.SkillSet
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

func (m *Mage) GetRotation() []*skills.AppliedSkill {
	return m.SkillSet.Rotation
}

func (m *Mage) AddSkillToRotation(boardName string, skillName string, fromX int, fromY int, toX int, toY int) {
	log.Debug("adding skill to rotation: %s: %d, %d -> %d, %d", skillName, fromX, fromY, toX, toY)
	m.SkillSet.Rotation = append(m.SkillSet.Rotation, &skills.AppliedSkill{
		boardName,
		1,
		fromX,
		fromY,
		toX,
		toY,
		skillName,
	},
	)
}

func (m *Mage) GetSkillSet() *skills.SkillSet {
	return m.SkillSet
}

func (m *Mage) SetSkillSet(ss *skills.SkillSet) {
	m.SkillSet = ss
}

func (m *Mage) LearnSkill(name string, skill skills.SkillFunc) {
	m.SkillSet.SkillBook[name] = skill
}

func (m *Mage) ApplySkills() {
	for _, app := range m.SkillSet.Rotation {
		if _, ok := m.SkillSet.SkillBook[app.Name]; !ok {
			log.Debug("no such skill in the book: %s", app.Name)
			continue
		}
		m.SkillSet.SkillBook[app.Name](app.Board, app.FromX, app.FromY, app.ToX, app.ToY)
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

// SetCoords sets unit x,y coords
func (m *Mage) SetCoords(X, Y int) {
	m.Figure.X = X
	m.Figure.Y = Y
}

// GetCoords gets unit x,y coords
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
