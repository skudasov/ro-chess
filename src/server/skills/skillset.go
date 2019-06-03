package skills

// SkillFunc skill mechanics goes here
type SkillFunc func(string, int, int, int, int)

// Entity of skill application from one unit to another
type AppliedSkill struct {
	// How many times skill will be applied
	Board string
	Times int
	FromX int
	FromY int
	ToX   int
	ToY   int
	Name  string
}

type SkillSet struct {
	// SKill rotation to be applied in order
	Rotation []*AppliedSkill
	// All learned skills mechanics will be here
	// front-end doesn't need to know them
	SkillBook map[string]SkillFunc `json:"-"`
}

func NewEmptySkillSet() *SkillSet {
	return &SkillSet{
		Rotation:  make([]*AppliedSkill, 0),
		SkillBook: map[string]SkillFunc{},
	}
}
