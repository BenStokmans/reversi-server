package handlers

import (
	"github.com/BenStokmans/reversi-server/game"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"
)

func handlePlayMove(msg *anypb.Any, client *game.Client) error {
	playMove := &game.PlayMove{}
	err := anypb.UnmarshalTo(msg, playMove, proto.UnmarshalOptions{})
	if err != nil {
		return err
	}

	if client.State.Game == nil {
		return nil
	}

	resp := &game.PlayMoveResponse{
		Success: true,
	}
	err = client.State.Game.Move(client, playMove)
	if err != nil {
		resp.Success = false
		reason := err.Error()
		resp.Error = &reason
	}

	client.Send(resp)
	return nil
}
