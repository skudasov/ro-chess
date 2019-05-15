package internal

import (
	"github.com/f4hrenh9it/ro-chess/src/server/conf"
	"github.com/f4hrenh9it/ro-chess/src/server/game"
	"github.com/f4hrenh9it/ro-chess/src/server/msg"
	"github.com/name5566/leaf/gate"
)

// Module represents game module
type Module struct {
	*gate.Gate
}

// OnInit adds default gate object
func (m *Module) OnInit() {
	m.Gate = &gate.Gate{
		MaxConnNum:      conf.Server.MaxConnNum,
		PendingWriteNum: conf.PendingWriteNum,
		MaxMsgLen:       conf.MaxMsgLen,
		WSAddr:          conf.Server.WSAddr,
		HTTPTimeout:     conf.HTTPTimeout,
		CertFile:        conf.Server.CertFile,
		KeyFile:         conf.Server.KeyFile,
		TCPAddr:         conf.Server.TCPAddr,
		LenMsgLen:       conf.LenMsgLen,
		LittleEndian:    conf.LittleEndian,
		Processor:       msg.Processor,
		AgentChanRPC:    game.ChanRPC,
	}
}
