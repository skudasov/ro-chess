package entity

type figure struct {
	Type       string
	Name       string
	VisualMark string
	Owner      string
	Movable    bool
	Active     bool
	PrevX      int
	PrevY      int
	X          int
	Y          int
	Alive      bool
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
