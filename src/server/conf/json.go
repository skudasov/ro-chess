package conf

import (
	"encoding/json"
	"github.com/name5566/leaf/log"
	"io/ioutil"
	"os/exec"
	"strings"
)

// Server config struct
var Server struct {
	LogLevel    string
	LogPath     string
	WSAddr      string
	CertFile    string
	KeyFile     string
	TCPAddr     string
	MaxConnNum  int
	ConsolePort int
	ProfilePath string

	PlayerHP       int
	PlayerMP       int
	FigurePoolType string
	BoardSizeX     int
	BoardSizeY     int

	ZoneScoreYTop              int
	ZoneStartYTopPlayerFrom    int
	ZoneStartYTopPlayerTo      int
	ZoneStartYBottomPlayerFrom int
	ZoneStartYBottomPlayerTo   int
	ZoneScoreYBottom           int
}

func rootPath() (string, error) {
	path, err := exec.Command("git", "rev-parse", "--show-toplevel").Output()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(path)), nil
}

func init() {
	rp, err := rootPath()
	if err != nil {
		log.Fatal("%v", err)
	}
	data, err := ioutil.ReadFile(rp + "/bin/conf/server.json")
	if err != nil {
		log.Fatal("%v", err)
	}
	err = json.Unmarshal(data, &Server)
	if err != nil {
		log.Fatal("%v", err)
	}
}
