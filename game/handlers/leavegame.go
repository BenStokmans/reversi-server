package handlers

import (
	"github.com/BenStokmans/reversi-server/game"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"
)

func handleLeaveGame(msg *anypb.Any, client *game.Client) error {
	leaveGame := &game.LeaveGame{}
	err := anypb.UnmarshalTo(msg, leaveGame, proto.UnmarshalOptions{})
	if err != nil {
		return err
	}

	if client.State.Game == nil {
		return nil
	}
	client.State.Game.RemovePlayer(client)
	client.State.Game = nil

	return nil
}
