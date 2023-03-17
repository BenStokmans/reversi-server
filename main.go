//go:generate protoc -I=../ --go_out=. reversi-server/reversi.proto
package main

import (
	"github.com/BenStokmans/reversi-server/game"
	"github.com/BenStokmans/reversi-server/game/handlers"
	"github.com/sirupsen/logrus"
)

func main() {
	logrus.SetLevel(logrus.DebugLevel)
	server := game.NewServer(handlers.HandleMessage)
	server.Start()
}
