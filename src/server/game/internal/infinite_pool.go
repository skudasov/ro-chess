package internal

import (
	"fmt"
	e "github.com/f4hrenh9it/ro-chess/src/server/entity"
)

// InfiniteFigurePool used for test purpose, Board testing
type InfiniteFigurePool struct {
	PoolLvl int
	Race    race
	Figures []e.Figurable
}

// NewInfiniteFigurePool creates new pool with units for testing
func NewInfiniteFigurePool(race race, poolLvl int) *InfiniteFigurePool {
	return &InfiniteFigurePool{
		poolLvl,
		race,
		make([]e.Figurable, 0),
	}
}

// visualize visualizes units in pool
func (m *InfiniteFigurePool) Visualize() {
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
func (m *InfiniteFigurePool) Fill(owner string) {
	//fill it from race units pool/locker
	m.Figures = append(m.Figures,
		e.NewWarrior(e.OptOwner(owner)),
		e.NewMage(e.OptOwner(owner)),
		e.NewWarrior(e.OptOwner(owner)),
		e.NewWarrior(e.OptOwner(owner)),
		e.NewWarrior(e.OptOwner(owner)),
	)
}

// Get gets copy of unit from pool
func (m *InfiniteFigurePool) Get(x int) e.Figurable {
	return m.Figures[x].Clone()
}

// GetFigures gets all figures from pool
func (m *InfiniteFigurePool) GetFigures() []e.Figurable {
	return m.Figures
}
