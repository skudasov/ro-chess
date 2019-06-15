package entity

// CombatEvent represents combat event for front-ent animation
type CombatEvent struct {
	// think about consumption and different skills animation
	From *Point
	To   *Point
	Dmg  int
	// if it's empty, figure uses auto-attack
	AnimationName string
}
