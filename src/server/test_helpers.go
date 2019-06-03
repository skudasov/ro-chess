package server

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"github.com/f4hrenh9it/ro-chess/src/server/conf"
	"github.com/f4hrenh9it/ro-chess/src/server/msg"
	"log"
	"net"
)

type client struct {
	Name string
	net.Conn
}

type cJoin struct {
	Join msg.Join
}

type cJoined struct {
	Joined msg.Joined
}

type cGameStarted struct {
	GameStarted msg.GameStarted
}

type cTurnFigurePool struct {
	TurnFigurePool msg.TurnFigurePool
}

type cMoveFigure struct {
	MoveFigure msg.MoveFigure
}

type cGameError struct {
	GameError msg.GameError
}

type cDisconnect struct {
	Disconnect msg.Disconnect
}

type cEndTurn struct {
	EndTurn msg.EndTurn
}

type cYourTurn struct {
	YourTurn msg.YourTurn
}

type cActivateFigure struct {
	ActivateFigure msg.ActivateFigure
}

type cUpdateBatch struct {
	UpdateBatch msg.UpdateBatch
}

type cYouWin struct {
	YouWin msg.YouWin
}

type cYouLose struct {
	YouLose msg.YouLose
}

type cCastSkill struct {
	CastSkill msg.CastSkill
}

func oneUnitBoard() {
	conf.Server.BoardSizeX = 1
	conf.Server.BoardSizeY = 3
	conf.Server.ZoneStartYTopPlayerFrom = 1
	conf.Server.ZoneStartYTopPlayerTo = 1
	conf.Server.ZoneStartYBottomPlayerFrom = 1
	conf.Server.ZoneStartYBottomPlayerTo = 1
	conf.Server.ZoneScoreYTop = 0
	conf.Server.ZoneScoreYBottom = 2
}

func twoUnitsCombatBoard() {
	conf.Server.FigurePoolType = "constdmg"
	conf.Server.BoardSizeX = 1
	conf.Server.BoardSizeY = 4
	conf.Server.ZoneStartYTopPlayerFrom = 1
	conf.Server.ZoneStartYTopPlayerTo = 1
	conf.Server.ZoneStartYBottomPlayerFrom = 2
	conf.Server.ZoneStartYBottomPlayerTo = 2
	conf.Server.ZoneScoreYTop = 0
	conf.Server.ZoneScoreYBottom = 3
}

func newClient(name string) *client {
	p1, err := net.Dial("tcp", "127.0.0.1:3563")
	if err != nil {
		log.Fatal(err)
	}
	return &client{name, p1}
}

func unpack(resBytes []byte, reply string) interface{} {
	var unpacker map[string]interface{}
	if err := json.Unmarshal(resBytes, &unpacker); err != nil {
		log.Fatal(err)
	}
	// unpack msg name, the only key
	for k := range unpacker {
		if k != reply {
			fmt.Printf("unwanted msg received of type: %s\n, wants: %s\n", k, reply)
			return nil
		}
	}
	switch reply {
	case "Joined":
		var res cJoined
		if err := json.Unmarshal(resBytes, &res); err != nil {
			log.Fatal(err)
		}
		return res
	case "GameError":
		var res cGameError
		if err := json.Unmarshal(resBytes, &res); err != nil {
			log.Fatal(err)
		}
		return res
	case "MoveFigure":
		var res cMoveFigure
		if err := json.Unmarshal(resBytes, &res); err != nil {
			log.Fatal(err)
		}
		return res
	case "Disconnect":
		var res cDisconnect
		if err := json.Unmarshal(resBytes, &res); err != nil {
			log.Fatal(err)
		}
		return res
	case "GameStarted":
		var res cGameStarted
		if err := json.Unmarshal(resBytes, &res); err != nil {
			log.Fatal(err)
		}
		return res
	case "TurnFigurePool":
		var res cTurnFigurePool
		if err := json.Unmarshal(resBytes, &res); err != nil {
			log.Fatal(err)
		}
		return res
	case "EndTurn":
		var res cEndTurn
		if err := json.Unmarshal(resBytes, &res); err != nil {
			log.Fatal(err)
		}
		return res
	case "YourTurn":
		var res cYourTurn
		if err := json.Unmarshal(resBytes, &res); err != nil {
			log.Fatal(err)
		}
		return res
	case "ActivateFigure":
		var res cActivateFigure
		if err := json.Unmarshal(resBytes, &res); err != nil {
			log.Fatal(err)
		}
		return res
	case "UpdateBatch":
		var res cUpdateBatch
		if err := json.Unmarshal(resBytes, &res); err != nil {
			log.Fatal(err)
		}
		return res
	case "YouWin":
		var res cYouWin
		if err := json.Unmarshal(resBytes, &res); err != nil {
			log.Fatal(err)
		}
		return res
	case "YouLose":
		var res cYouLose
		if err := json.Unmarshal(resBytes, &res); err != nil {
			log.Fatal(err)
		}
		return res
	default:
		log.Fatal("unknown msg received")
		return nil
	}
}

func (m *client) send(msg interface{}) {
	data, _ := json.Marshal(msg)
	msgData := make([]byte, 2+len(data))
	binary.BigEndian.PutUint16(msgData, uint16(len(data)))
	copy(msgData[2:], data)
	fmt.Printf("%s sending msg: %s\n", m.Name, msgData)
	m.Conn.Write(msgData)
}

func (m *client) sendAndRead(msg interface{}, reply string) interface{} {
	data, _ := json.Marshal(msg)
	msgData := make([]byte, 2+len(data))
	binary.BigEndian.PutUint16(msgData, uint16(len(data)))
	copy(msgData[2:], data)
	fmt.Printf("%s sending msg: %s\n", m.Name, msgData)
	m.Conn.Write(msgData)

	resBytes := readTCP(m.Conn)
	fmt.Printf("%s receiving msg: %s\n", m.Name, resBytes)
	res := unpack(resBytes, reply)
	return res
}

func (m *client) read(reply string) interface{} {
	resBytes := readTCP(m.Conn)
	fmt.Printf("%s receiving msg: %s\n", m.Name, resBytes)
	res := unpack(resBytes, reply)
	return res
}

func readTCP(conn net.Conn) []byte {
	msgLenBytes := make([]byte, 2)
	if _, err := conn.Read(msgLenBytes); err != nil {
		log.Fatal(err)
	}
	msgLen := binary.BigEndian.Uint16(msgLenBytes)
	fmt.Printf("msg received: %d\n", msgLen)
	answer := make([]byte, msgLen)
	if _, err := conn.Read(answer); err != nil {
		log.Fatal(err)
	}
	return answer
}
