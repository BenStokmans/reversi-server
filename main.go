//go:generate protoc -I=../ --go_out=. reversi-server/reversi.proto
package main

import "github.com/BenStokmans/reversi-server/game"

func main() {
	server := game.Server{}
	server.Start()
}
