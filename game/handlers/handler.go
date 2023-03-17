package handlers

import (
	"github.com/BenStokmans/reversi-server/game"
	"google.golang.org/protobuf/types/known/anypb"
)

type messageHandler func(msg *anypb.Any, client *game.Client) error

var handlers = map[string]messageHandler{
	"reversi.CreateGame": handleCreateGame,
	"reversi.JoinGame":   handleJoinGame,
	"reversi.PlayMove":   handlePlayMove,
	"reversi.LeaveGame":  handleLeaveGame,
}

func HandleMessage(msg *anypb.Any, client *game.Client) error {
	handler, ok := handlers[string(msg.MessageName())]
	if !ok {
		return nil
	}
	return handler(msg, client)
}
