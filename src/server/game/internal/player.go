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
