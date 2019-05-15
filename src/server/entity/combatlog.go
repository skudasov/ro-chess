package entity

// CombatEvent represents combat event for front-ent animation
type CombatEvent struct {
	X, Y int
	Dmg  int
	Crit bool
}
