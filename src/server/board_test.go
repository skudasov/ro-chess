package server

import (
	"fmt"
	"github.com/f4hrenh9it/ro-chess/src/server/conf"
	e "github.com/f4hrenh9it/ro-chess/src/server/entity"
	"github.com/f4hrenh9it/ro-chess/src/server/game"
	"github.com/f4hrenh9it/ro-chess/src/server/gate"
	"github.com/f4hrenh9it/ro-chess/src/server/login"
	"github.com/f4hrenh9it/ro-chess/src/server/msg"
	"github.com/name5566/leaf"
	lconf "github.com/name5566/leaf/conf"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

const (
	board = "abc"
)

func runLeaf() {
	lconf.LogLevel = conf.Server.LogLevel
	lconf.LogPath = conf.Server.LogPath
	lconf.LogFlag = conf.LogFlag
	lconf.ConsolePort = conf.Server.ConsolePort
	lconf.ProfilePath = conf.Server.ProfilePath
	conf.Server.FixedTurns = "p2"

	go leaf.Run(
		game.Module,
		gate.Module,
		login.Module,
	)
	time.Sleep(1 * time.Second)
}

func TestBoardSuite(t *testing.T) {
	runLeaf()
	t.Run("TestConnect", func(t *testing.T) {
		p1 := newClient("p1")
		defer p1.Conn.Close()
		res := p1.sendAndRead(
			cJoin{msg.Join{"p1", "p1"}},
			"Joined")
		fmt.Printf("board is: %s\n", res)
		res2 := p1.sendAndRead(cDisconnect{msg.Disconnect{"p1", board, "leaving"}}, "Disconnect")
		fmt.Printf("disconnect res: %s", res2)
	})

	t.Run("TestBoardRotation", func(t *testing.T) {
		g := newBoardGame(t, board)
		g.end()
		g = newBoardGame(t, board)
		g.end()
	})

	t.Run("TestMoveNoBoard", func(t *testing.T) {
		p1 := newClient("p1")
		defer p1.Conn.Close()
		res := p1.sendAndRead(
			cMoveFigure{msg.MoveFigure{"p1", board, 0, 0, 1}},
			"GameError",
		)
		assert.Equal(t, cGameError{msg.GameError{"no board"}}, res)
	})

	t.Run("TestFigureInCell", func(t *testing.T) {
		g := newBoardGame(t, board)
		defer g.end()

		p2movMsg := cMoveFigure{msg.MoveFigure{"p2", g.board, 0, 0, 14}}
		g.P2.sendAndRead(p2movMsg, "UpdateBatch")
		g.P1.read("UpdateBatch")
		//move again
		res3 := g.P2.sendAndRead(p2movMsg, "GameError")
		assert.Equal(t, cGameError{msg.GameError{"cell not empty"}}, res3)
	})

	t.Run("TestMoveFromPool", func(t *testing.T) {
		g := newBoardGame(t, board)
		defer g.end()

		p2movMsg := cMoveFigure{msg.MoveFigure{"p2", g.board, 0, 0, 14}}
		p2UpdMsg := cUpdateBatch{
			msg.UpdateBatch{
				Figures: []e.Figurable{
					e.NewWarrior(
						e.OptOwner("p2"),
						e.OptMovable(false),
						e.OptX(0),
						e.OptY(14),
					),
				},
			},
		}
		res2 := g.P2.sendAndRead(p2movMsg, "UpdateBatch")
		assert.Equal(t, p2UpdMsg, res2)
		res1 := g.P1.read("UpdateBatch")
		assert.Equal(t, p2UpdMsg, res1)
	})

	t.Run("TestApplyAttackCombo", func(t *testing.T) {
		g := newBoardGame(t, board)
		defer g.end()

		g.p2MovesFromPool(1, 0, 12)
		g.p2MovesFromPool(1, 0, 13)
		resp1, resp2 := g.p2MovesFromPool(1, 0, 14)
		// attack buff of +1-+3 applied
		updMsg := cUpdateBatch{
			msg.UpdateBatch{
				Figures: []e.Figurable{
					e.NewMage(
						e.OptOwner("p2"),
						e.OptMovable(false),
						e.OptX(0),
						e.OptY(12),
						e.OptAttackMin(8),
						e.OptAttackMax(13),
					),
					e.NewMage(
						e.OptOwner("p2"),
						e.OptMovable(false),
						e.OptX(0),
						e.OptY(13),
						e.OptAttackMin(8),
						e.OptAttackMax(13),
					),
					e.NewMage(
						e.OptOwner("p2"),
						e.OptMovable(false),
						e.OptX(0),
						e.OptY(14),
						e.OptAttackMin(8),
						e.OptAttackMax(13),
					),
				},
			},
		}
		assert.Equal(t, updMsg, resp1)
		assert.Equal(t, updMsg, resp2)
	})

	t.Run("TestApplyDefenseCombo", func(t *testing.T) {
		g := newBoardGame(t, board)
		defer g.end()

		g.p2MovesFromPool(1, 0, 14)
		g.p2MovesFromPool(1, 1, 14)
		resp1, resp2 := g.p2MovesFromPool(1, 2, 14)
		// attack buff of +1-+3 applied
		updMsg := cUpdateBatch{
			msg.UpdateBatch{
				Figures: []e.Figurable{
					e.NewMage(
						e.OptOwner("p2"),
						e.OptMovable(false),
						e.OptX(0),
						e.OptY(14),
						e.OptDefence(2),
					),
					e.NewMage(
						e.OptOwner("p2"),
						e.OptMovable(false),
						e.OptX(1),
						e.OptY(14),
						e.OptDefence(2),
					),
					e.NewMage(
						e.OptOwner("p2"),
						e.OptMovable(false),
						e.OptX(2),
						e.OptY(14),
						e.OptDefence(2),
					),
				},
			},
		}
		assert.Equal(t, updMsg, resp1)
		assert.Equal(t, updMsg, resp2)
	})

	t.Run("TestEndTurn", func(t *testing.T) {
		g := newBoardGame(t, board)
		defer g.end()

		g.p2MovesFromPool(0, 0, 14)

		// p2 ending turn
		g.P2.send(cEndTurn{msg.EndTurn{"p2", board}})
		g.P1.send(cEndTurn{msg.EndTurn{"p1", board}})
		g.P2.read("UpdateBatch")
		g.P1.read("UpdateBatch")
		yourTurn := g.P1.read("YourTurn")
		assert.Equal(t, cYourTurn{msg.YourTurn{}}, yourTurn)
		yourTurn2 := g.P2.read("YourTurn")
		assert.Equal(t, cYourTurn{msg.YourTurn{}}, yourTurn2)
	})

	t.Run("TestOutOfStartZone", func(t *testing.T) {
		g := newBoardGame(t, board)
		defer g.end()

		// p2 zones
		p2movMsg := cMoveFigure{msg.MoveFigure{"p2", board, 0, 10, 14}}
		outErr := g.P2.sendAndRead(p2movMsg, "GameError")
		assert.Equal(t, cGameError{msg.GameError{"figure out of start zone"}}, outErr)

		p2movMsg2 := cMoveFigure{msg.MoveFigure{"p2", board, 0, 8, 10}}
		outErr2 := g.P2.sendAndRead(p2movMsg2, "GameError")
		assert.Equal(t, cGameError{msg.GameError{"figure out of start zone"}}, outErr2)

		g.turnEnds()

		// p1 y zone
		p1movMsg := cMoveFigure{msg.MoveFigure{"p1", board, 0, 0, 5}}
		outErr3 := g.P1.sendAndRead(p1movMsg, "GameError")
		assert.Equal(t, cGameError{msg.GameError{"figure out of start zone"}}, outErr3)
	})

	t.Run("TestActivateFigure", func(t *testing.T) {
		g := newBoardGame(t, board)
		defer g.end()

		g.p2MovesFromPool(0, 8, 14)

		activateMsg := cActivateFigure{
			msg.ActivateFigure{
				"p2", board, 8, 14, true,
			},
		}
		resWalking := g.P2.sendAndRead(activateMsg, "ActivateFigure")
		assert.Equal(t, activateMsg, resWalking)
		resWalking2 := g.P1.read("ActivateFigure")
		assert.Equal(t, activateMsg, resWalking2)
	})

	t.Run("TestActivateFigureNoFigure", func(t *testing.T) {
		g := newBoardGame(t, board)
		defer g.end()

		g.p2MovesFromPool(0, 7, 14)

		activateMsg := cActivateFigure{
			msg.ActivateFigure{
				"p2", board, 8, 14, true,
			},
		}
		errMsg := cGameError{msg.GameError{"no figure in the cell"}}
		resWalking := g.P2.sendAndRead(activateMsg, "GameError")
		assert.Equal(t, errMsg, resWalking)
	})

	t.Run("TestFigureMovingAfterActiveSet", func(t *testing.T) {
		g := newBoardGame(t, board)
		defer g.end()

		// p2 moving
		g.p2MovesFromPool(0, 8, 14)
		g.p2ActivatesFigure(8, 14)
		resp1, _ := g.turnEnds()

		// p2 ending turn asserting updates
		peonUpdateMsg := cUpdateBatch{
			msg.UpdateBatch{
				Figures: []e.Figurable{
					e.NewWarrior(
						e.OptOwner("p2"),
						e.OptMovable(false),
						e.OptActive(true),
						e.OptPrevX(8),
						e.OptPrevY(14),
						e.OptX(8),
						e.OptY(13),
					),
				},
			},
		}
		assert.Equal(t, peonUpdateMsg, resp1)
	})

	t.Run("TestP1TakesDamage", func(t *testing.T) {
		oneUnitBoard()
		g := newBoardGame(t, board)
		defer g.end()

		g.p2MovesFromPool(0, 0, 1)
		g.p2ActivatesFigure(0, 1)
		g.turnEnds()
		upd1, upd2 := g.turnEnds()
		updMsg := cUpdateBatch{
			msg.UpdateBatch{
				Players: []e.Player{{"p1", 97, 0}},
				Figures: []e.Figurable{
					e.NewWarrior(
						e.OptOwner("p2"),
						e.OptMovable(false),
						e.OptAlive(false),
						e.OptActive(true),
						e.OptPrevX(0),
						e.OptPrevY(1),
						e.OptX(0),
						e.OptY(0),
					),
				},
			},
		}
		assert.Equal(t, updMsg, upd1)
		assert.Equal(t, updMsg, upd2)
	})

	t.Run("TestP2TakesDamage", func(t *testing.T) {
		oneUnitBoard()
		g := newBoardGame(t, board)
		defer g.end()

		g.p1MovesFromPool(0, 0, 1)
		g.p1ActivatesFigure(0, 1)
		g.turnEnds()
		upd1, upd2 := g.turnEnds()
		updMsg := cUpdateBatch{
			msg.UpdateBatch{
				Players: []e.Player{{"p2", 97, 0}},
				Figures: []e.Figurable{
					e.NewWarrior(
						e.OptOwner("p1"),
						e.OptMovable(false),
						e.OptActive(true),
						e.OptAlive(false),
						e.OptPrevX(0),
						e.OptPrevY(1),
						e.OptX(0),
						e.OptY(2),
					),
				},
			},
		}
		assert.Equal(t, updMsg, upd1)
		assert.Equal(t, updMsg, upd2)
	})

	t.Run("TestP1Loses", func(t *testing.T) {
		conf.Server.PlayerHP = 1
		oneUnitBoard()
		g := newBoardGame(t, board)
		defer g.end()

		// p2 moving
		g.p2MovesFromPool(0, 0, 1)
		g.p2ActivatesFigure(0, 1)
		g.P1.send(cEndTurn{msg.EndTurn{"p1", board}})
		g.P2.send(cEndTurn{msg.EndTurn{"p2", board}})
		g.P1.read("UpdateBatch")
		g.P2.read("UpdateBatch")
		g.P1.read("YouLose")
		g.P2.read("YouWin")
	})

	t.Run("TestP2Loses", func(t *testing.T) {
		conf.Server.PlayerHP = 1
		oneUnitBoard()
		g := newBoardGame(t, board)
		defer g.end()

		g.p1MovesFromPool(0, 0, 1)
		g.p1ActivatesFigure(0, 1)
		g.P2.send(cEndTurn{msg.EndTurn{"p2", board}})
		g.P1.send(cEndTurn{msg.EndTurn{"p1", board}})
		g.P2.read("UpdateBatch")
		g.P1.read("UpdateBatch")
		g.P1.read("YouWin")
		g.P2.read("YouLose")
	})

	t.Run("TestP1FigureDiesInCombat", func(t *testing.T) {
		twoUnitsCombatBoard()
		g := newBoardGame(t, board)
		defer g.end()

		g.p2MovesFromPool(0, 0, 2)
		g.p2ActivatesFigure(0, 2)
		g.p1MovesFromPool(0, 0, 1)
		g.p1ActivatesFigure(0, 1)
		updMsg := cUpdateBatch{
			msg.UpdateBatch{
				CombatLog: []e.CombatEvent{
					{X: 0, Y: 1, Dmg: -5, Crit: false},
					{X: 0, Y: 2, Dmg: -5, Crit: false},
					{X: 0, Y: 1, Dmg: -5, Crit: false},
				},
				Figures: []e.Figurable{
					e.NewConstDmgWarrior(
						e.OptOwner("p1"),
						e.OptMovable(false),
						e.OptActive(true),
						e.OptAlive(false),
						e.OptPrevX(0),
						e.OptPrevY(0),
						e.OptX(0),
						e.OptY(1),
						e.OptHP(0),
					),
					e.NewConstDmgWarrior(
						e.OptOwner("p2"),
						e.OptMovable(false),
						e.OptActive(true),
						e.OptPrevX(0),
						e.OptPrevY(2),
						e.OptX(0),
						e.OptY(1),
						e.OptHP(5),
					),
				},
			},
		}
		upd1, upd2 := g.turnEnds()
		assert.Equal(t, updMsg, upd1)
		assert.Equal(t, updMsg, upd2)
	})

	t.Run("TestP2FigureDiesInCombat", func(t *testing.T) {
		// TODO: test is usless now, refactor initiative so p2 can win
		twoUnitsCombatBoard()
		g := newBoardGame(t, board)
		defer g.end()

		g.p1MovesFromPool(0, 0, 1)
		g.p1ActivatesFigure(0, 1)
		g.p2MovesFromPool(0, 0, 2)
		g.p2ActivatesFigure(0, 2)
		upd1, upd2 := g.turnEnds()
		updMsg := cUpdateBatch{
			msg.UpdateBatch{
				CombatLog: []e.CombatEvent{
					{X: 0, Y: 2, Dmg: -5, Crit: false},
					{X: 0, Y: 1, Dmg: -5, Crit: false},
					{X: 0, Y: 2, Dmg: -5, Crit: false},
				},
				Figures: []e.Figurable{
					e.NewConstDmgWarrior(
						e.OptOwner("p2"),
						e.OptMovable(false),
						e.OptActive(true),
						e.OptAlive(false),
						e.OptPrevX(0),
						e.OptPrevY(0),
						e.OptX(0),
						e.OptY(2),
						e.OptHP(0),
					),
					e.NewConstDmgWarrior(
						e.OptOwner("p1"),
						e.OptMovable(false),
						e.OptActive(true),
						e.OptPrevX(0),
						e.OptPrevY(1),
						e.OptX(0),
						e.OptY(2),
						e.OptHP(5),
					),
				},
			},
		}
		assert.Equal(t, updMsg, upd1)
		assert.Equal(t, updMsg, upd2)
	})

	t.Run("TestSkillApplication", func(t *testing.T) {
		twoUnitsCombatBoard()
		g := newBoardGame(t, board)
		defer g.end()

		g.p1MovesFromPool(0, 0, 1)
		g.p1ActivatesFigure(0, 1)
		g.p2MovesFromPool(0, 0, 2)
		g.p2ActivatesFigure(0, 2)
		skillMsg := cCastSkill{msg.CastSkill{"p2", board, 0, 1, 0, 2, "fireball"}}
		g.P2.send(skillMsg)
		g.turnEnds()
	})
}
