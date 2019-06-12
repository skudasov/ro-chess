package server

import (
	"github.com/f4hrenh9it/ro-chess/src/server/msg"
	"github.com/stretchr/testify/assert"
	"testing"
)

type boardGame struct {
	t     *testing.T
	P1    *client
	P2    *client
	board string
}

func newBoardGame(t *testing.T, board string) *boardGame {
	p1 := newClient("p1")
	p2 := newClient("p2")
	p1.sendAndRead(cJoin{msg.Join{"p1", "p1"}}, "Joined")
	started := p2.sendAndRead(cJoin{msg.Join{"p2", "p2"}}, "GameStarted")
	assert.Equal(t, cGameStarted{msg.GameStarted{"p2"}}, started)
	p1.read("GameStarted")
	p2.read("TurnFigurePool")
	p2.read("YourTurn")
	return &boardGame{t, p1, p2, board}
}

func (m *boardGame) p2MovesFromPool(PoolX, X, Y int) (interface{}, interface{}) {
	p2movMsg := cMoveFigure{msg.MoveFigure{"p2", m.board, PoolX, X, Y}}
	resp2 := m.P2.sendAndRead(p2movMsg, "UpdateBatch")
	resp1 := m.P1.read("UpdateBatch")
	return resp1, resp2
}

func (m *boardGame) p1MovesFromPool(PoolX, X, Y int) (interface{}, interface{}) {
	p1movMsg := cMoveFigure{msg.MoveFigure{"p1", m.board, PoolX, X, Y}}
	resp1 := m.P1.sendAndRead(p1movMsg, "MoveFigure")
	resp2 := m.P2.read("MoveFigure")
	return resp1, resp2
}

func (m *boardGame) turnEnds() (interface{}, interface{}, interface{}, interface{}) {
	m.P2.send(cEndTurn{msg.EndTurn{"p2", m.board}})
	m.P1.send(cEndTurn{msg.EndTurn{"p1", m.board}})
	sfResp2 := m.P2.read("UpdateBatch")
	sfResp1 := m.P1.read("UpdateBatch")
	resp2 := m.P2.read("UpdateBatch")
	resp1 := m.P1.read("UpdateBatch")
	m.P2.read("YourTurn")
	m.P1.read("YourTurn")
	return sfResp1, sfResp2, resp1, resp2
}

func (m *boardGame) p2EndsTurn() {
	m.P2.send(cEndTurn{msg.EndTurn{"p2", m.board}})
	m.P2.read("UpdateBatch")
	m.P1.read("UpdateBatch")
	m.P2.read("TurnEnded")
	m.P1.read("YourTurn")
}

func (m *boardGame) p1EndsTurn() {
	m.P1.send(cEndTurn{msg.EndTurn{"p1", m.board}})
	m.P1.read("UpdateBatch")
	m.P2.read("UpdateBatch")
	m.P1.read("TurnEnded")
	m.P2.read("YourTurn")
}

func (m *boardGame) p2ActivatesFigure(X, Y int) {
	activateMsg := cActivateFigure{msg.ActivateFigure{"p2", m.board, X, Y, true}}
	m.P2.sendAndRead(activateMsg, "ActivateFigure")
	m.P1.read("ActivateFigure")
}

func (m *boardGame) p1ActivatesFigure(X, Y int) {
	activateMsg := cActivateFigure{msg.ActivateFigure{"p1", m.board, X, Y, true}}
	m.P1.sendAndRead(activateMsg, "ActivateFigure")
	m.P2.read("ActivateFigure")
}

func (m *boardGame) end() {
	m.P1.sendAndRead(cDisconnect{msg.Disconnect{"p1", m.board, "leaving"}}, "Disconnect")
	m.P2.sendAndRead(cDisconnect{msg.Disconnect{"p2", m.board, "leaving"}}, "Disconnect")
}
