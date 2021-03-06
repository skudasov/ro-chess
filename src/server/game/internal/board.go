package internal

import (
	"fmt"
	"github.com/f4hrenh9it/ro-chess/src/server/conf"
	e "github.com/f4hrenh9it/ro-chess/src/server/entity"
	"github.com/f4hrenh9it/ro-chess/src/server/msg"
	"github.com/name5566/leaf/log"
	"math/rand"
	"sort"
	"time"
)

// no mutexes is needed, read leaf docs about handler goroutines
// e.g. read func (c *LinearContext) Go(f func(), cb func())

// BS represents all boards, every Board is a game
var BS = make(map[string]*Board)

// SkillLibrary All skills mechanics goes here
type SkillLibrary map[string]e.SkillFunc

// SL global skill library to learn from
// every skill has access to all board figures and can update figures, players, create combat log
// If skill has *from* pair and *to* is nil - it's buff skill
// or target of application will be found by filtering units on board (player reference required?)
// If skill has both pairs, second pair will be target of skill or center of AOE skill
// after players/figures transformation one must attach CombatEvents in order:
// 1. MP/HP consumption of source (caster)
// 2. If it's targeted skill or we need animation (ex. firebolt) send
// TODO: think about skills as general transformation function that contains two slices of PointsOfAction which consists of
// TODO: 1. X,Y
// TODO: 2. Consumption that will be rendered (-HP, +MP, etc.)
// TODO: AnimationName
// TODO: So clog for fireball will be:
// TODO: 1. clog with animation of casting, consumption of MP and XY of caster
// TODO: 2. clog with animation of flying fireball, two XY, from and to
// TODO: 3. clog with units taking damage, where from = nil, and to is a slice of diminished dmg
// TODO: can we simplify this?
var SL = SkillLibrary{
	"firebolt": func(boardName string, from e.Point, to e.Point, uf *[]e.Figurable, up *[]e.Player, clog *[]e.CombatEvent) {
		fromFigure := BS[boardName].Canvas[from.Y][from.X].Figure
		toFigure := BS[boardName].Canvas[to.Y][to.X].Figure
		log.Debug("interacting figures: %s -> %s", fromFigure.GetName(), toFigure.GetName())
		log.Debug("casting %s: %d, %d -> %d, %d", "firebolt", from.X, from.Y, to.X, to.Y)
		toFigure.SetHP(toFigure.GetHP() - 200)
		*uf = append(*uf, toFigure)
		*clog = append(*clog, e.CombatEvent{})
	},
}

var visTemplate = "[ %2s %2s %2d %2d %5s %2d ]"

func createBoard(boardName string, players map[string]*player, xSize, ySize int) (*Board, string) {
	log.Debug("creating Board %s", boardName)
	c := make([][]*cell, ySize)
	for i := range c {
		c[i] = make([]*cell, xSize)
		for j := range c[i] {
			c[i][j] = &cell{}
		}
	}

	b := &Board{
		Players:  players,
		Canvas:   c,
		TurnEnds: make(map[string]bool),
	}
	b.createStartZones()
	for _, p := range players {
		p.Board = b
	}
	BS[boardName] = b
	log.Debug("Board created\n")
	return BS[boardName], ""
}

func firstTurn(players map[string]*player) string {
	s := make([]string, 0)
	for k := range players {
		s = append(s, k)
	}
	rand.Seed(time.Now().Unix())
	return s[rand.Intn(len(s))]
}

func (m *Board) setActive(X, Y int, status bool) error {
	if m.Canvas[Y][X].Figure == nil {
		return newErrGameNoFigureInCell()
	}
	m.Canvas[Y][X].Figure.Activate(status)
	return nil
}

func (m *Board) createStartZones() {
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

func (m *Board) checkStartZone(x, y int, side side) bool {
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

func (m *Board) moveFromPool(pToken string, side side, poolX, X, Y int) (e.Figurable, error) {
	player := m.Players[pToken]
	figure := player.FigurePool.Get(poolX)
	if !m.checkStartZone(X, Y, side) {
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
	m.Figures = append(m.Figures, figure)
	m.VisualizeAll()
	return figure, nil
}

func (m *Board) visualize() {
	log.Debug("Board visualization")
	for i := range m.Canvas {
		for j := range m.Canvas[i] {
			if m.Canvas[i][j].Figure != nil {
				f := m.Canvas[i][j].Figure
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
				fmt.Printf(visTemplate, "", "", 0, 0, "-", 0)
			}
		}
		fmt.Print("\n")
	}
}

func (m *Board) broadcast(message interface{}) {
	sortedPKeys := make([]string, 0)
	for pkey := range m.Players {
		sortedPKeys = append(sortedPKeys, pkey)
	}
	sort.Strings(sortedPKeys)
	for k := range m.Players {
		m.Players[k].Agent.WriteMsg(message)
	}
}

func (m *Board) endGame(loser *string) {
	player := m.Players[*loser]
	player.Opponent.Agent.WriteMsg(&msg.YouWin{})
	player.Agent.WriteMsg(&msg.YouLose{})
}
