package gate

import (
	"github.com/f4hrenh9it/ro-chess/src/server/game"
	"github.com/f4hrenh9it/ro-chess/src/server/msg"
)

func init() {
	msg.Processor.SetRouter(&msg.Join{}, game.ChanRPC)
	msg.Processor.SetRouter(&msg.Disconnect{}, game.ChanRPC)
	msg.Processor.SetRouter(&msg.EndTurn{}, game.ChanRPC)
	msg.Processor.SetRouter(&msg.MoveFigure{}, game.ChanRPC)
	msg.Processor.SetRouter(&msg.ActivateFigure{}, game.ChanRPC)
	msg.Processor.SetRouter(&msg.CastSkill{}, game.ChanRPC)
	msg.Processor.SetRouter(&msg.LearnSkill{}, game.ChanRPC)
}
