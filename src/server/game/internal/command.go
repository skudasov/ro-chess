package internal

import (
	"github.com/name5566/leaf/log"
)

// visualize board and pool data
func (m *board) VisualizeAll() {
	m.Visualize()
	for t, p := range m.Players {
		log.Debug("%s pool visualization", t)
		p.FigurePool.Visualize()
	}
}
