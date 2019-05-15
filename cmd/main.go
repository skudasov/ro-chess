package main

import (
	"github.com/f4hrenh9it/ro-chess/src/server/conf"
	"github.com/f4hrenh9it/ro-chess/src/server/game"
	"github.com/f4hrenh9it/ro-chess/src/server/gate"
	"github.com/f4hrenh9it/ro-chess/src/server/login"
	"github.com/name5566/leaf"
	lconf "github.com/name5566/leaf/conf"
)

func main() {
	lconf.LogLevel = conf.Server.LogLevel
	lconf.LogPath = conf.Server.LogPath
	lconf.LogFlag = conf.LogFlag
	lconf.ConsolePort = conf.Server.ConsolePort
	lconf.ProfilePath = conf.Server.ProfilePath

	leaf.Run(
		game.Module,
		gate.Module,
		login.Module,
	)
}
