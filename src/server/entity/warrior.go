package entity

import (
	"github.com/name5566/leaf/log"
)

// Warrior unit
type Warrior struct {
	figure
	SkillSet *SkillSet
}

// NewWarrior creates new peon unit
func NewWarrior(opts ...func(m *figure)) *Warrior {
	p := &Warrior{figure{
		Type:       "warrior",
		Name:       "Warrior",
		VisualMark: "W",
		Movable:    true,
		Active:     false,
		PrevX:      0,
		PrevY:      0,
		Alive:      true,
		HP:         10,
		MP:         10,
		AttackMin:  1,
		AttackMax:  3,
		Armor:      0,
		Initiative: 1,
	}, nil}
	for _, o := range opts {
		o(&p.figure)
	}
	return p
}

// NewConstDmgWarrior creates new peon unit with const dps
func NewConstDmgWarrior(opts ...func(m *figure)) *Warrior {
	p := &Warrior{figure{
		Type:       "warrior",
		Name:       "Warrior",
		VisualMark: "W",
		Movable:    true,
		Active:     false,
		PrevX:      0,
		PrevY:      0,
		Alive:      true,
		HP:         10,
		MP:         10,
		AttackMin:  5,
		AttackMax:  5,
		Armor:      0,
		Initiative: 1,
	}, nil}
	for _, o := range opts {
		o(&p.figure)
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

// AddSkillToRotation adds skill to skill rotation
func (m *Warrior) AddSkillToRotation(boardName string, skillName string, from Point, to Point) {
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

// Clone copies unit
func (m *Warrior) Clone() Figurable {
	cloned := *m
	return &cloned
}
