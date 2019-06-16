package msg

import (
	"encoding/json"
	"fmt"
	e "github.com/f4hrenh9it/ro-chess/src/server/entity"
	ljson "github.com/name5566/leaf/network/json"
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
	Processor.Register(&YourTurn{})
	Processor.Register(&LvlUp{})
	Processor.Register(&LearnSkill{})
	Processor.Register(&CastSkill{})
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
	Figures []e.Figurable
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
	Figure e.Figurable
}

// LvlUp msg
type LvlUp struct {
	AvailableSkills []string
}

// LearnSkill msg (after player receives LvlUp msg he can learn a skill)
type LearnSkill struct {
	Token string
	Board string
	Name  string
	From  e.Point
}

// UpdateBatch msg
type UpdateBatch struct {
	Players   []e.Player      `json:"Players,omitempty"`
	CombatLog []e.CombatEvent `json:"CombatLog,omitempty"`
	Figures   []e.Figurable   `json:"Figures,omitempty"`
}

// CastSkill msg
type CastSkill struct {
	Token string
	Board string
	From  e.Point
	To    e.Point
	Name  string
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
	ce.Figures = make([]e.Figurable, len(rawMsgsFigures))

	var m map[string]interface{}
	for index, rawMessage := range rawMsgsFigures {
		err = json.Unmarshal(*rawMessage, &m)
		if err != nil {
			return err
		}

		switch m["Type"] {
		case "warrior":
			var p e.Warrior
			if err := json.Unmarshal(*rawMessage, &p); err != nil {
				return err
			}
			ce.Figures[index] = &p
		case "mage":
			var g e.Mage
			if err := json.Unmarshal(*rawMessage, &g); err != nil {
				return err
			}
			ce.Figures[index] = &g
		default:
			return fmt.Errorf("unsupported type found when unmarshalling: %s", m["Type"])
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
		ce.Figures = make([]e.Figurable, len(figures))
	}

	var players []e.Player
	if _, ok := objMap["Players"]; ok {
		if err = json.Unmarshal(*objMap["Players"], &players); err != nil {
			return err
		}
		ce.Players = players
	}

	var combatLog []e.CombatEvent
	if _, ok := objMap["CombatLog"]; ok {
		if err = json.Unmarshal(*objMap["CombatLog"], &combatLog); err != nil {
			return err
		}
		ce.CombatLog = combatLog
	}

	var m map[string]interface{}
	for index, rawMessage := range figures {
		err = json.Unmarshal(*rawMessage, &m)
		if err != nil {
			return err
		}

		switch m["Type"] {
		case "warrior":
			var p e.Warrior
			if err := json.Unmarshal(*rawMessage, &p); err != nil {
				return err
			}
			ce.Figures[index] = &p
		case "mage":
			var g e.Mage
			if err := json.Unmarshal(*rawMessage, &g); err != nil {
				return err
			}
			ce.Figures[index] = &g
		default:
			return fmt.Errorf("unsupported type found when unmarshalling: %s", m["Type"])
		}
	}
	return nil
}
