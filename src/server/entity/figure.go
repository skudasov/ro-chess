package entity

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
