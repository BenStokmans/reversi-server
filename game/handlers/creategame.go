package handlers

import (
	"github.com/BenStokmans/reversi-server/game"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"
)

func HandleCreateGame(msg *anypb.Any) error {
	createGame := &game.CreateGame{}
	err := anypb.UnmarshalTo(msg, createGame, proto.UnmarshalOptions{})
	if err != nil {
		return err
	}
	return nil
}
