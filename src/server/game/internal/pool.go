package internal

import (
	"fmt"
	e "github.com/f4hrenh9it/ro-chess/src/server/entity"
)

// FigurePool contains figures that user allowed to draw in that turn
type FigurePool struct {
	PoolLvl int
	Race    race
	Figures []e.Figurable
}

// NewFigurePool creates new figure pool
func NewFigurePool(race race, poolLvl int) *FigurePool {
	return &FigurePool{
		poolLvl,
		race,
		make([]e.Figurable, 0),
	}
}

// Visualize visualizes units in pool
func (m *FigurePool) Visualize() {
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
func (m *FigurePool) Fill(owner string) {
	//fill it from race units pool/locker
	m.Figures = append(m.Figures,
		e.NewWarrior(e.OptOwner(owner)),
		e.NewMage(e.OptOwner(owner)),
		e.NewWarrior(e.OptOwner(owner)),
		e.NewWarrior(e.OptOwner(owner)),
		e.NewWarrior(e.OptOwner(owner)),
	)
}

// Get gets unit from pool
func (m *FigurePool) Get(x int) e.Figurable {
	return m.Figures[x]
}

// GetFigures gets all figures from pool
func (m *FigurePool) GetFigures() []e.Figurable {
	return m.Figures
}
