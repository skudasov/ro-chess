package internal

import (
	"github.com/f4hrenh9it/ro-chess/src/server/conf"
	"github.com/f4hrenh9it/ro-chess/src/server/entity"
	"github.com/f4hrenh9it/ro-chess/src/server/msg"
	"github.com/name5566/leaf/gate"
	"github.com/name5566/leaf/log"
	"reflect"
)

type registering struct {
	Token string
	Name  string
	Agent gate.Agent
}

var registeringQueue = make([]*registering, 0)

func handler(m interface{}, h interface{}) {
	skeleton.RegisterChanRPC(reflect.TypeOf(m), h)
}

func init() {
	handler(&msg.Join{}, handleJoin)
	handler(&msg.Disconnect{}, handleDisconnect)
	handler(&msg.EndTurn{}, handleEndTurn)
	handler(&msg.MoveFigure{}, handleMoveFromPool)
	handler(&msg.CastSkill{}, handleCastSkill)
	handler(&msg.ActivateFigure{}, handleActivateFigure)
}

func handleDisconnect(args []interface{}) {
	m := args[0].(*msg.Disconnect)
	a := args[1].(gate.Agent)
	log.Debug("disconnect msg: %v", m.Token)
	for i, rq := range registeringQueue {
		if rq.Token == m.Token {
			log.Debug("deleting token: %s", m.Token)
			registeringQueue = append(registeringQueue[:i], registeringQueue[i+1:]...)
		}
	}
	board := BS[m.Board]

	if board == nil {
		a.WriteMsg(&msg.Disconnect{m.Token, "abc", "ok"})
		return
	}
	for token := range board.Players {
		if token == m.Token {
			delete(board.Players, token)
			// TODO: game over here and p1 wins
			a.WriteMsg(&msg.Disconnect{m.Token, "abc", "ok"})
		}
	}
	log.Debug("Board players left: %d", len(board.Players))
	if len(board.Players) == 0 {
		log.Debug("deleting Board %s", m.Board)
		delete(BS, m.Board)
	}
	for _, r := range registeringQueue {
		log.Debug("queue after disconnect: token: %s, name: %s", r.Token, r.Name)
	}
}

func handleJoin(args []interface{}) {
	m := args[0].(*msg.Join)
	a := args[1].(gate.Agent)

	// validate token somewhere
	log.Debug("join msg: %v", m.Token)
	registeringQueue = append(registeringQueue, &registering{m.Token, m.Name, a})

	for _, r := range registeringQueue {
		log.Debug("token: %s, name: %s", r.Token, r.Name)
	}

	// make auth here or in login module
	// pull info from db after login

	// fetch player info here, create player and start game

	// make match making here, for now just start if we have two
	bname := "abc"
	if len(registeringQueue) == 2 {
		players := initPlayers()

		log.Debug("sizes: x = %d, y = %d", conf.Server.BoardSizeX, conf.Server.BoardSizeY)
		board, turn := createBoard(
			bname,
			players,
			conf.Server.BoardSizeX, conf.Server.BoardSizeY,
		)
		board.Broadcast(&msg.GameStarted{turn})
		// send pool figures for the first turn
		pName := board.Turn
		player := board.Players[pName]

		poolFigures := player.FigurePool.GetFigures()

		player.Agent.WriteMsg(&msg.TurnFigurePool{poolFigures})
		player.Agent.WriteMsg(&msg.YourTurn{})
		// here we are waiting for client to end turn by himself
		return
	}

	a.WriteMsg(&msg.Joined{
		Board: bname,
	})
}

func handleMoveFromPool(args []interface{}) {
	m := args[0].(*msg.MoveFigure)
	a := args[1].(gate.Agent)
	log.Debug("move msg")
	if _, ok := BS[m.Board]; !ok {
		a.WriteMsg(&msg.GameError{newErrGameNoBoard().Error()})
		return
	}
	board := BS[m.Board]
	player := board.Players[m.Token]

	updatedFigures := make([]entity.Figurable, 0)
	f, err := board.MoveFromPool(
		m.Token,
		player.Side,
		m.PoolX, m.ToX, m.ToY,
	)
	if err != nil {
		a.WriteMsg(&msg.GameError{err.Error()})
		return
	}
	updatedFigures = append(updatedFigures, f)

	player.putCombo(f)

	updatedByCombos, err := player.applyCombo()
	if err != nil {
		a.WriteMsg(&msg.GameError{err.Error()})
	}
	if len(updatedByCombos) != 0 {
		board.Broadcast(&msg.UpdateBatch{nil, nil, updatedByCombos})
		return
	}
	log.Debug("figures here: %s", updatedFigures)
	board.Broadcast(&msg.UpdateBatch{nil, nil, updatedFigures})
}

func handleEndTurn(args []interface{}) {
	m := args[0].(*msg.EndTurn)
	a := args[1].(gate.Agent)
	log.Debug("end turn msg")
	if _, ok := BS[m.Board]; !ok {
		a.WriteMsg(&msg.GameError{newErrGameNoBoard().Error()})
		return
	}
	board := BS[m.Board]
	// process turn mechanics here
	log.Debug("turn ended by %s", m.Token)
	board.TurnEnds[m.Token] = true
	if len(board.TurnEnds) != 2 {
		log.Debug("not enough players for attack phase, waiting")
		return
	}

	log.Debug("turn phase: Battle")
	updatedFigures := make([]entity.Figurable, 0)
	updatedPlayers := make([]entity.Player, 0)
	combatLogs := make([]entity.CombatEvent, 0)

	deadPlayerName := board.processCombat(&updatedFigures, &updatedPlayers, &combatLogs)
	if deadPlayerName != nil {
		player := board.Players[*deadPlayerName]
		player.Opponent.Agent.WriteMsg(&msg.YouWin{})
		player.Agent.WriteMsg(&msg.YouLose{})
		return
	}
	board.Broadcast(&msg.UpdateBatch{updatedPlayers, combatLogs, updatedFigures})
	board.Broadcast(&msg.YourTurn{})
	log.Debug("turn phase: Planning")
	board.TurnEnds = make(map[string]bool)
}

func handleCastSkill(args []interface{}) {
	m := args[0].(*msg.CastSkill)
	a := args[1].(gate.Agent)
	log.Debug("cast skill msg")
	if _, ok := BS[m.Board]; !ok {
		a.WriteMsg(&msg.GameError{newErrGameNoBoard().Error()})
		return
	}
	board := BS[m.Board]
	fromUnit := board.Canvas[m.FromY][m.FromX]
	fromUnit.Figure.LearnSkill("fireball", SL["fireball"])
	fromUnit.Figure.AddSkillToRotation(m.Board, m.Name, m.FromX, m.FromY, m.ToX, m.ToY)
}

func handleActivateFigure(args []interface{}) {
	m := args[0].(*msg.ActivateFigure)
	a := args[1].(gate.Agent)
	log.Debug("figure activate msg")
	if _, ok := BS[m.Board]; !ok {
		a.WriteMsg(&msg.GameError{newErrGameNoBoard().Error()})
		return
	}
	board := BS[m.Board]
	// if it's not a squad limit amount of moving figures so you cannot get all area revealed
	// no one wants to die alone, without honor!
	err := board.SetActive(m.X, m.Y, m.Active)
	if err != nil {
		a.WriteMsg(&msg.GameError{err.Error()})
		return
	}
	log.Debug("sending activate msg")
	board.Broadcast(m)
}
