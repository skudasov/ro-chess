package internal

import (
	"fmt"
	"github.com/f4hrenh9it/ro-chess/src/server/conf"
	"github.com/f4hrenh9it/ro-chess/src/server/entity"
	"github.com/name5566/leaf/log"
	"math/rand"
	"time"
)

// no mutexes is needed, read leaf docs about handler goroutines
// e.g. read func (c *LinearContext) Go(f func(), cb func())
var bs = make(map[string]*board)

var visTemplate = "[ %2s %2s %2d %2d %5s %2d ]"

func createBoard(boardName string, players map[string]*player, xSize, ySize int) (*board, string) {
	log.Debug("creating board %s", boardName)
	c := make([][]*cell, ySize)
	for i := range c {
		c[i] = make([]*cell, xSize)
		for j := range c[i] {
			c[i][j] = &cell{}
		}
	}
	var ft string
	if conf.Server.FixedTurns != "" {
		ft = conf.Server.FixedTurns
	} else {
		ft = firstTurn(players)
	}

	b := &board{
		Players: players,
		Canvas:  c,
		Turn:    ft}
	b.CreateStartZones()
	for _, p := range players {
		p.Board = b
	}
	bs[boardName] = b
	log.Debug("board created\n")
	return bs[boardName], ft
}

func firstTurn(players map[string]*player) string {
	s := make([]string, 0)
	for k := range players {
		s = append(s, k)
	}
	rand.Seed(time.Now().Unix())
	return s[rand.Intn(len(s))]
}

func (m *board) SetActive(X, Y int, status bool) error {
	if m.Canvas[Y][X].Figure == nil {
		return newErrGameNoFigureInCell()
	}
	m.Canvas[Y][X].Figure.Activate(status)
	return nil
}

func (m *board) CreateStartZones() {
	log.Debug("creating start zones")
	for i := 0; i < conf.Server.BoardSizeX; i++ {
		m.ZoneStartTopY = append(m.ZoneStartTopY, i)
	}
	for i := conf.Server.ZoneStartYTopPlayerFrom; i <= conf.Server.ZoneStartYTopPlayerTo; i++ {
		m.ZoneStartTopX = append(m.ZoneStartTopX, i)
	}
	for i := 0; i < conf.Server.BoardSizeX; i++ {
		m.ZoneStartBottomY = append(m.ZoneStartBottomY, i)
	}
	for i := conf.Server.ZoneStartYBottomPlayerFrom; i <= conf.Server.ZoneStartYBottomPlayerTo; i++ {
		m.ZoneStartBottomX = append(m.ZoneStartBottomX, i)
	}
	log.Debug("top zone y: %d\n", m.ZoneStartTopX)
	log.Debug("top zone x: %d\n", m.ZoneStartTopY)
	log.Debug("bottom zone y: %d\n", m.ZoneStartBottomX)
	log.Debug("bottom zone x: %d\n", m.ZoneStartBottomY)
}

func (m *board) CheckStartZone(x, y int, side side) bool {
	log.Debug("x, y = %d, %d", x, y)
	if side == top {
		log.Debug("top")
		for _, i := range m.ZoneStartTopX {
			for _, j := range m.ZoneStartTopY {
				if x == j && y == i {
					return true
				}
			}
		}
	} else {
		log.Debug("bottom")
		for _, i := range m.ZoneStartBottomX {
			for _, j := range m.ZoneStartBottomY {
				if x == j && y == i {
					return true
				}
			}
		}
	}
	return false
}

func (m *board) MoveFromPool(pToken string, side side, poolX, X, Y int) (entity.Figurable, error) {
	player := m.Players[pToken]
	figure := player.FigurePool.Get(poolX)
	if !m.CheckStartZone(X, Y, side) {
		return nil, newErrFigureMoveOutOfStartZone()
	}
	if figure == nil {
		return nil, newErrGameNoFigureInCell()
	}
	if m.Canvas[Y][X].Figure != nil {
		return nil, newErrCellNotEmpty()
	}
	if !figure.GetMovable() {
		return nil, newErrGameFigureNotMovable()
	}
	targetCell := m.Canvas[Y][X]

	// for now after first move figure is locked forever
	log.Debug("moving %s from pool %d, to %d, %d\n",
		figure.GetName(), poolX, X, Y)
	targetCell.Figure = figure
	targetCell.Figure.SetMovable(false)
	targetCell.Figure.SetCoords(X, Y)
	player.BoardFigures = append(player.BoardFigures, figure)
	m.VisualizeAll()
	return figure, nil
}

func (m *board) Visualize() {
	log.Debug("board visualization")
	for i := range m.Canvas {
		for j := range m.Canvas[i] {
			if m.Canvas[i][j].Figure != nil {
				f := m.Canvas[i][j].Figure
				fmt.Printf(
					visTemplate,
					f.GetOwner(),
					f.GetVisualMark(),
					f.GetHP(),
					f.GetMP(),
					f.GetAttackStr(),
					f.GetDefence(),
				)
			} else {
				fmt.Printf(visTemplate, "", "", 0, 0, "-", 0)
			}
		}
		fmt.Print("\n")
	}
}

func (m *board) Broadcast(message interface{}) {
	for _, p := range m.Players {
		p.Agent.WriteMsg(message)
	}
}
