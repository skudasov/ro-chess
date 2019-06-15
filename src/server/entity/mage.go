package entity

import (
	"github.com/name5566/leaf/log"
)

// Mage unit
type Mage struct {
	figure
	SkillSet *SkillSet
}

// NewMage creates new grunt unit
func NewMage(opts ...func(m *figure)) *Mage {
	g := &Mage{figure{}, nil}
	g.Type = "mage"
	g.Name = "Mage"
	g.VisualMark = "M"
	g.Movable = true
	g.Active = false
	g.PrevX = 0
	g.PrevY = 0
	g.Alive = true
	g.HP = 30
	g.MP = 0
	g.AttackMin = 7
	g.AttackMax = 10
	g.Armor = 1
	for _, o := range opts {
		o(&g.figure)
	}
	return g
}

// GetRotation Gets skills rotation
func (m *Mage) GetRotation() []*AppliedSkill {
	return m.SkillSet.Rotation
}

// AddSkillToRotation Adds skill to rotation
func (m *Mage) AddSkillToRotation(boardName string, skillName string, from Point, to Point) {
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

// Clone copies unit
func (m *Mage) Clone() Figurable {
	cloned := *m
	return &cloned
}
