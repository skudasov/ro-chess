package internal

import (
	"github.com/f4hrenh9it/ro-chess/src/server/base"
	"github.com/name5566/leaf/module"
)

var (
	skeleton = base.NewSkeleton()
	// ChanRPC var
	ChanRPC = skeleton.ChanRPCServer
)

// Module represents game module
type Module struct {
	*module.Skeleton
}

// OnInit puts skeleton to module
func (m *Module) OnInit() {
	m.Skeleton = skeleton
}

// OnDestroy does nothing
func (m *Module) OnDestroy() {

}
