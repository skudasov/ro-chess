package internal

import (
	"github.com/f4hrenh9it/ro-chess/src/server/conf"
	"github.com/f4hrenh9it/ro-chess/src/server/entity"
	"github.com/name5566/leaf/log"
)

const (
	maxComboStacks = 5
)

// Creates players object, fill unit pools, set opponent
func initPlayers() map[string]*player {

	players := make(map[string]*player, 0)

	var next bool
	for _, reg := range registeringQueue {
		var side side
		var opponentScoreZone int
		if !next {
			side = top
			opponentScoreZone = conf.Server.ZoneScoreYBottom
		} else {
			side = bottom
			opponentScoreZone = conf.Server.ZoneScoreYTop
		}
		var fp entity.Poolable
		switch conf.Server.FigurePoolType {
		case "infinite":
			fp = NewInfiniteFigurePool(orcs, 0)
		case "constdmg":
			fp = NewConstDmgFigurePool(orcs, 0)
		default:
			fp = NewFigurePool(orcs, 0)
		}
		fp.Fill(reg.Name)

		p := &player{
			reg.Name,
			reg.Agent,
			nil,
			reg.Name,
			side,
			nil,
			conf.Server.PlayerHP,
			conf.Server.PlayerMP,
			// opposite scoring zones
			opponentScoreZone,
			comboStore{},
			make([]entity.Figurable, 0),
			fp,
		}
		players[reg.Name] = p
		next = true
	}
	for _, p := range players {
		log.Debug("player %s has start zone %s and scoring at %d row", p.Name, p.Side, p.ZoneScoreY)
	}
	players[registeringQueue[0].Name].Opponent = players[registeringQueue[1].Name]
	players[registeringQueue[1].Name].Opponent = players[registeringQueue[0].Name]
	return players
}

// Move all player figures one cell forward
func (m *player) moveFigures(updated *[]entity.Figurable) {
	log.Debug("moving figures for %s", m.Info)
	for _, newFigure := range m.BoardFigures {
		// all active figures moves 1 cell up or down, depends on player side
		if newFigure.GetActive() == true {
			prevTurnX, prevTurnY := newFigure.GetCoords()
			// set prev coords so we can update all using one figure data
			newFigure.SetPrevCoords(prevTurnX, prevTurnY)
			var newX, newY int
			if m.Side == top {
				newX, newY = prevTurnX, prevTurnY+1
			} else {
				newX, newY = prevTurnX, prevTurnY-1
			}
			prevCell := m.Board.Canvas[prevTurnY][prevTurnX]

			targetCell := m.Board.Canvas[newY][newX]

			// all scoring and battle mechanics happens before this moment
			// see damagePlayer, attackFigures
			targetCell.Figure = newFigure
			prevCell.Figure = nil
			newFigure.SetCoords(newX, newY)
			m.Board.VisualizeAll()

			*updated = append(*updated, newFigure)
		}
	}
}

func (m *player) deleteBoardFigures(bfs *[]entity.Figurable) {
	for i, bf := range *bfs {
		if bf.GetAlive() == false {
			*bfs = append((*bfs)[:i], (*bfs)[i+1:]...)
		}
	}
}

// Damages opponent if your figures reach opponent scoring zone
func (m *player) damagePlayer(players *[]entity.Player, figures *[]entity.Figurable) bool {
	// if our figure comes to opponent dmg zone, deal dmg equal to attack
	log.Debug("scoring figures for %s", m.Info)
	scored := make([]int, 0)
	for i, f := range m.BoardFigures {
		if x, y := f.GetCoords(); y == m.ZoneScoreY {
			// think about damage model for players, is min-max from PerformAttack() is ok?
			magicNumber := 3
			log.Debug("%s takes damage from figure is score zone!", m.Opponent.Name)
			log.Debug("figure: %s", f.GetName())
			log.Debug("figure damage: %d", magicNumber)
			f.SetAlive(false)
			m.Board.Canvas[y][x] = nil
			m.Opponent.HP = m.Opponent.HP - magicNumber
			*players = append(*players, entity.Player{
				Name: m.Opponent.Name,
				HP:   m.Opponent.HP,
				MP:   m.Opponent.MP,
			},
			)
			scored = append(scored, i)
			f.SetAlive(false)
			*figures = append(*figures, f)
			// if it's squad, do the math and replace all squad units in one time
		}
	}
	m.deleteBoardFigures(&m.BoardFigures)
	if m.Opponent.HP <= 0 {
		return true
	}
	return false
}

// Creating combat for every figure that faces opponent figure in that turn, fills combat log with events.
// After fighting is done, deletes killed figures from your BoardFigures or opponent
func (m *player) attackFigures(figures *[]entity.Figurable, clog *[]entity.CombatEvent) {
	for _, f := range m.BoardFigures {
		var enemyX, enemyY int
		atkX, atkY := f.GetCoords()
		if m.Side == top {
			enemyX, enemyY = atkX, atkY+1
		} else {
			enemyX, enemyY = atkX, atkY-1
		}
		attacker := f
		defender := m.Board.Canvas[enemyY][enemyX].Figure
		if defender != nil {
			log.Debug("conflict!")
			defX, defY := defender.GetCoords()
			for {
				aDmg := attacker.PerformAttack()
				log.Debug("attacker dmg: %d", aDmg)
				defender.SetHP(defender.GetHP() - aDmg)
				*clog = append(*clog, entity.CombatEvent{defX, defY, -aDmg, false})
				if defender.GetHP() <= 0 {
					defender.SetAlive(false)
					m.Board.Canvas[defY][defX].Figure = nil
					log.Debug("defender on %d, %d, figure of %s dies", atkX, atkY, defender.GetOwner())
					*figures = append(*figures, defender)
					break
				}
				dDmg := defender.PerformAttack()
				log.Debug("defender dmg: %d", dDmg)
				attacker.SetHP(attacker.GetHP() - dDmg)
				*clog = append(*clog, entity.CombatEvent{atkX, atkY, -dDmg, false})
				if attacker.GetHP() <= 0 {
					attacker.SetAlive(false)
					m.Board.Canvas[atkY][atkX].Figure = nil
					log.Debug("attacker on %d, %d, figure of %s dies", atkX, atkY, attacker.GetOwner())
					*figures = append(*figures, attacker)
					break
				}
			}
		}
	}
	m.deleteBoardFigures(&m.BoardFigures)
	m.deleteBoardFigures(&m.Opponent.BoardFigures)
}

// putCombo adds figure to comboStore for x and y dimensions
func (m *player) putCombo(f entity.Figurable) {
	x, y := f.GetCoords()
	kX := comboKey{"x", f.GetName(), x}
	kY := comboKey{"y", f.GetName(), y}

	if m.ComboStore[kX] == nil {
		m.ComboStore[kX] = &comboVal{}
	}
	m.ComboStore[kX].UnitCoord = append(m.ComboStore[kX].UnitCoord, y)
	if m.ComboStore[kY] == nil {
		m.ComboStore[kY] = &comboVal{}
	}
	m.ComboStore[kY].UnitCoord = append(m.ComboStore[kY].UnitCoord, x)
	log.Debug("combos: %s", m.ComboStore)
}

// applyCombo applies buff on player groups of units
// every player has
func (m *player) applyCombo() ([]entity.Figurable, error) {
	log.Debug("searching combos for %s", m.Name)
	updated := make([]entity.Figurable, 0)
	for k, combo := range m.ComboStore {
		if len(combo.UnitCoord) == 3 {
			combo.Stacks++
			if combo.Stacks >= maxComboStacks {
				return nil, nil
			}
			// TODO: refactor all Figurable slices to sets! crutchy!
			log.Debug("combo found of 3 %s", k.Class)
			var f entity.Figurable
			var keyDim = k.Indent
			var comboType comboType
			for _, comboDim := range combo.UnitCoord {
				switch k.Axis {
				case "x":
					f = m.Board.Canvas[comboDim][keyDim].Figure
					comboType = attack3Combo
				case "y":
					f = m.Board.Canvas[keyDim][comboDim].Figure
					comboType = defense3Combo
				}
				if f == nil {
					log.Debug("f not found")
					return nil, newErrGameNoFigureInCell()
				}
				m.buffByFigureType(f, comboType)
				log.Debug("applying %s combo buff on f: %s, on %d, %d", comboType, f.GetName(), keyDim, comboDim)
				// TODO: switch by other characteristics here
				updated = append(updated, f)
			}
			log.Debug("%d stacks of combo applied", combo.Stacks)
		}
	}
	return updated, nil
}

func (m *player) buffByFigureType(f entity.Figurable, comboType comboType) {
	x, y := f.GetCoords()
	switch {
	case f.GetName() == "Grunt" && comboType == attack3Combo:
		log.Debug("%s on %d, %d receives minAtk: %d, maxAtk: %d", f.GetName(), x, y, 1, 3)
		f.AddAttack(1, 3)
	case f.GetName() == "Grunt" && comboType == defense3Combo:
		log.Debug("%s on %d, %d receives armor: %d", f.GetName(), x, y, 1)
		f.AddDefense(1)
	}
}
