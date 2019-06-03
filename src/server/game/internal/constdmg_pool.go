package internal

import (
	"fmt"
	e "github.com/f4hrenh9it/ro-chess/src/server/entity"
)

// ConstDmgFigurePool used for test purpose, combat testing
type ConstDmgFigurePool struct {
	PoolLvl int
	Race    race
	Figures []e.Figurable
}

// NewConstDmgFigurePool creates new pool with constant dmg units for testing
func NewConstDmgFigurePool(race race, poolLvl int) *ConstDmgFigurePool {
	return &ConstDmgFigurePool{
		poolLvl,
		race,
		make([]e.Figurable, 0),
	}
}

// Visualize visualizes units in pool
func (m *ConstDmgFigurePool) Visualize() {
	for i := range m.Figures {
		if m.Figures[i] != nil {
			f := m.Figures[i]
			fmt.Printf(
				visTemplate,
				f.GetOwnerName(),
				f.GetVisualMark(),
				f.GetHP(),
				f.GetMP(),
				f.GetAttackStr(),
				f.GetDefence(),
			)
		} else {
			fmt.Printf(visTemplate, "", "", 0, 0, 0, 0)
		}
		fmt.Print("\n")
	}
}

// Fill fills pool with units
func (m *ConstDmgFigurePool) Fill(owner string) {
	m.Figures = append(m.Figures,
		e.NewConstDmgWarrior(e.OptOwner(owner)),
		e.NewConstDmgWarrior(e.OptOwner(owner)),
		e.NewConstDmgWarrior(e.OptOwner(owner)),
		e.NewConstDmgWarrior(e.OptOwner(owner)),
		e.NewConstDmgWarrior(e.OptOwner(owner)),
	)
}

// Get gets copy of unit from pool
func (m *ConstDmgFigurePool) Get(x int) e.Figurable {
	return m.Figures[x].Clone()
}

// GetFigures gets all figures from pool
func (m *ConstDmgFigurePool) GetFigures() []e.Figurable {
	return m.Figures
}
