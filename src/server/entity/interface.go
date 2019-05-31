package entity

// Figurable represents actions with figures
type Figurable interface {
	Clonable
	Movable
	Buffable
	Combatable
	Visualizable
}

// Visualizable represents actions with visualization of figures
type Visualizable interface {
	GetName() string
	GetVisualMark() string
	GetOwnerName() string
}

// Movable represents actions to move figures on the board
type Movable interface {
	GetMovable() bool
	SetMovable(bool)
	Activate(bool)
	GetActive() bool
	SetCoords(int, int)
	GetCoords() (int, int)
	SetPrevCoords(int, int)
	GetPrevCoords() (int, int)
}

// Combatable represents actions performed in combat with other units
type Combatable interface {
	PerformAttack() int
	GetInitiative() int
	GetAlive() bool
	SetAlive(bool)
	SetHP(int)
	GetHP() int
	SetMP(int)
	GetMP() int
	GetDefence() int
	GetAttack() (int, int)
	GetAttackStr() string
}

// Clonable let object behind interface be cloned
type Clonable interface {
	Clone() Figurable
}

// Buffable represents actions for buff mechanics
type Buffable interface {
	AddAttack(int, int)
	AddDefense(int)
}

// Poolable represents actions with figures pool
type Poolable interface {
	GetFigures() []Figurable
	Fill(string)
	Get(int) Figurable
	Visualize()
}
