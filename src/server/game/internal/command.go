package internal

import (
	"github.com/name5566/leaf/log"
)

// VisualizeAll Visualizes Board and pool data
func (m *Board) VisualizeAll() {
	m.visualize()
	for t, p := range m.Players {
		log.Debug("%s pool visualization", t)
		p.FigurePool.Visualize()
	}
}
