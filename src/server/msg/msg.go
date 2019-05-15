package msg

import (
	"encoding/json"
	"github.com/f4hrenh9it/ro-chess/src/server/entity"
	ljson "github.com/name5566/leaf/network/json"
	"github.com/pkg/errors"
)

// Processor represents msg processor (json/pb)
var Processor = ljson.NewProcessor()

func init() {
	Processor.Register(&Join{})
	Processor.Register(&Joined{})
	Processor.Register(&GameStarted{})
	Processor.Register(&YouWin{})
	Processor.Register(&YouLose{})
	Processor.Register(&TurnFigurePool{})
	Processor.Register(&EndTurn{})
	Processor.Register(&TurnEnded{})
	Processor.Register(&YourTurn{})
	Processor.Register(&ActivateFigure{})
	Processor.Register(&FigureUpdate{})
	Processor.Register(&UpdateBatch{})
	Processor.Register(&Disconnect{})
	Processor.Register(&GameError{})
	Processor.Register(&MoveFigure{})
}

// Join msg
type Join struct {
	Token string
	Name  string
}

// Disconnect msg
type Disconnect struct {
	Token  string
	Board  string
	Reason string
}

// Joined msg
type Joined struct {
	Board string
}

// GameStarted msg
type GameStarted struct {
	Turn string
}

// TurnFigurePool msg
type TurnFigurePool struct {
	Figures []entity.Figurable
}

// GameError msg
type GameError struct {
	Msg string
}

// MoveFigure msg
type MoveFigure struct {
	Token string
	Board string
	PoolX int
	ToX   int
	ToY   int
}

// EndTurn msg
type EndTurn struct {
	Token string
	Board string
}

// YourTurn msg
type YourTurn struct{}

// TurnEnded msg
type TurnEnded struct{}

// ActivateFigure msg
type ActivateFigure struct {
	Token  string
	Board  string
	X      int
	Y      int
	Active bool
}

// FigureUpdate msg
type FigureUpdate struct {
	Figure entity.Figurable
}

// UpdateBatch msg
type UpdateBatch struct {
	Players   []entity.Player      `json:"Players,omitempty"`
	CombatLog []entity.CombatEvent `json:"CombatLog,omitempty"`
	Figures   []entity.Figurable   `json:"Figures,omitempty"`
}

// YouWin msg
type YouWin struct{}

// YouLose msg
type YouLose struct{}

// UnmarshalJSON allows to unmarshal types behind Figurable interface
func (ce *TurnFigurePool) UnmarshalJSON(b []byte) error {
	var objMap map[string]*json.RawMessage
	err := json.Unmarshal(b, &objMap)
	if err != nil {
		return err
	}
	var rawMsgsFigures []*json.RawMessage
	err = json.Unmarshal(*objMap["Figures"], &rawMsgsFigures)
	if err != nil {
		return err
	}
	ce.Figures = make([]entity.Figurable, len(rawMsgsFigures))

	var m map[string]map[string]interface{}
	for index, rawMessage := range rawMsgsFigures {
		err = json.Unmarshal(*rawMessage, &m)
		if err != nil {
			return err
		}

		switch m["Figure"]["Type"] {
		case "peon":
			var p entity.Peon
			if err := json.Unmarshal(*rawMessage, &p); err != nil {
				return err
			}
			ce.Figures[index] = &p
		case "grunt":
			var g entity.Grunt
			if err := json.Unmarshal(*rawMessage, &g); err != nil {
				return err
			}
			ce.Figures[index] = &g
		default:
			return errors.New("unsupported type found when unmarshalling")
		}
	}
	return nil
}

// UnmarshalJSON allows to unmarshal types behind Figurable interface
func (ce *UpdateBatch) UnmarshalJSON(b []byte) error {
	var objMap map[string]*json.RawMessage
	err := json.Unmarshal(b, &objMap)
	if err != nil {
		return err
	}
	var figures []*json.RawMessage
	if _, ok := objMap["Figures"]; ok {
		err = json.Unmarshal(*objMap["Figures"], &figures)
		if err != nil {
			return err
		}
		ce.Figures = make([]entity.Figurable, len(figures))
	}

	var players []entity.Player
	if _, ok := objMap["Players"]; ok {
		if err = json.Unmarshal(*objMap["Players"], &players); err != nil {
			return err
		}
		ce.Players = players
	}

	var combatLog []entity.CombatEvent
	if _, ok := objMap["CombatLog"]; ok {
		if err = json.Unmarshal(*objMap["CombatLog"], &combatLog); err != nil {
			return err
		}
		ce.CombatLog = combatLog
	}

	var m map[string]map[string]interface{}
	for index, rawMessage := range figures {
		err = json.Unmarshal(*rawMessage, &m)
		if err != nil {
			return err
		}

		switch m["Figure"]["Type"] {
		case "peon":
			var p entity.Peon
			if err := json.Unmarshal(*rawMessage, &p); err != nil {
				return err
			}
			ce.Figures[index] = &p
		case "grunt":
			var g entity.Grunt
			if err := json.Unmarshal(*rawMessage, &g); err != nil {
				return err
			}
			ce.Figures[index] = &g
		default:
			return errors.New("unsupported type found")
		}
	}
	return nil
}
