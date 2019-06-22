package internal

import (
	"github.com/f4hrenh9it/ro-chess/src/server/conf"
	"github.com/f4hrenh9it/ro-chess/src/server/entity"
	"github.com/f4hrenh9it/ro-chess/src/server/msg"
	"github.com/name5566/leaf/gate"
	"github.com/name5566/leaf/log"
	"reflect"
	"sort"
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
	handler(&msg.LearnSkill{}, handleLearnSkill)
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
		registeringQueue = make([]*registering, 0)
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

	bname := "abc"
	if len(registeringQueue) == 2 {
		players := initPlayers()

		log.Debug("sizes: x = %d, y = %d", conf.Server.BoardSizeX, conf.Server.BoardSizeY)
		board, _ := createBoard(
			bname,
			players,
			conf.Server.BoardSizeX, conf.Server.BoardSizeY,
		)
		// send pool figures for the first turn

		for _, p := range board.Players {
			poolFigures := p.FigurePool.GetFigures()
			p.Agent.WriteMsg(&msg.Joined{bname})
			p.Agent.WriteMsg(&msg.TurnFigurePool{poolFigures})
			p.Agent.WriteMsg(&msg.YourTurn{})
		}
		// here we are waiting for client to end turn by himself
		return
	}
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
	f, err := board.moveFromPool(
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
		board.broadcast(&msg.UpdateBatch{nil, nil, updatedByCombos})
		return
	}
	log.Debug("figures here: %s", updatedFigures)
	board.broadcast(&msg.UpdateBatch{nil, nil, updatedFigures})
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
	// sorting all figures by initiative
	sort.Sort(figuresByInitiative(board.Figures))

	updatedFigures := make([]entity.Figurable, 0)
	updatedPlayers := make([]entity.Player, 0)
	combatLogs := make([]entity.CombatEvent, 0)
	deadPlayerName := board.processSkillUpdates(&updatedFigures, &updatedPlayers, &combatLogs)
	board.broadcast(&msg.UpdateBatch{updatedPlayers, combatLogs, updatedFigures})
	if deadPlayerName != nil {
		board.endGame(deadPlayerName)
		return
	}

	updatedFigures = make([]entity.Figurable, 0)
	updatedPlayers = make([]entity.Player, 0)
	combatLogs = make([]entity.CombatEvent, 0)
	deadPlayerName = board.processCombatUpdates(&updatedFigures, &updatedPlayers, &combatLogs)
	board.broadcast(&msg.UpdateBatch{updatedPlayers, combatLogs, updatedFigures})
	if deadPlayerName != nil {
		board.endGame(deadPlayerName)
		return
	}

	board.broadcast(&msg.YourTurn{})
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
	fromUnit := board.Canvas[m.From.Y][m.From.X]
	from := entity.Point{m.From.X, m.From.Y}
	to := entity.Point{m.To.X, m.To.Y}
	// TODO: check for skill availability here?
	fromUnit.Figure.AddSkillToRotation(m.Board, m.Name, from, to)
}

func handleLearnSkill(args []interface{}) {
	m := args[0].(*msg.LearnSkill)
	a := args[1].(gate.Agent)
	log.Debug("learn skill msg")
	if _, ok := BS[m.Board]; !ok {
		a.WriteMsg(&msg.GameError{newErrGameNoBoard().Error()})
		return
	}
	board := BS[m.Board]
	fromUnit := board.Canvas[m.From.Y][m.From.X]
	fromUnit.Figure.LearnSkill("firebolt", SL["firebolt"])
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
	err := board.setActive(m.X, m.Y, m.Active)
	if err != nil {
		a.WriteMsg(&msg.GameError{err.Error()})
		return
	}
	log.Debug("sending activate msg")
	board.broadcast(m)
}
