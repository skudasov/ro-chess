package entity

// SkillFunc skill mechanics goes here
type SkillFunc func(string, Pair, Pair, *[]Figurable, *[]Player, *[]CombatEvent)

// AppliedSkill Entity of skill application from one unit to another
type AppliedSkill struct {
	Board string
	// How many times skill will be applied
	Times int
	From  Pair
	To    Pair
	Name  string
}

// SkillSet Entity of current skills figure may have
type SkillSet struct {
	// SKill rotation to be applied in order
	Rotation []*AppliedSkill
	// All learned skills mechanics will be here
	// front-end doesn't need to know them
	SkillBook map[string]SkillFunc `json:"-"`
}

// NewEmptySkillSet Creates new empty skillset
func NewEmptySkillSet() *SkillSet {
	return &SkillSet{
		Rotation:  make([]*AppliedSkill, 0),
		SkillBook: map[string]SkillFunc{},
	}
}
