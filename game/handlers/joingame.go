package handlers

import (
	"github.com/BenStokmans/reversi-server/game"
	"github.com/BenStokmans/reversi-server/snowflake"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"
	"time"
)

func handleJoinGame(msg *anypb.Any, client *game.Client) error {
	joinGame := &game.JoinGame{}
	err := anypb.UnmarshalTo(msg, joinGame, proto.UnmarshalOptions{})
	if err != nil {
		return err
	}

	resp := &game.JoinGameResponse{}
	if client.State.Game != nil {
		errMsg := "already in a game"
		resp.Error = &errMsg
		client.Send(resp)
		return nil
	}

	color, err := game.ClientJoin(client, snowflake.Snowflake(joinGame.GameId))
	if err != nil {
		errMsg := err.Error()
		resp.Error = &errMsg
		client.Send(resp)
		return nil
	}
	resp.Color = color
	client.Send(resp)

	// wait for a short time to allow the player to join process to complete
	time.Sleep(5 * time.Millisecond)

	err = client.State.Game.Start()
	if err != nil {
		return err
	}

	return nil
}
