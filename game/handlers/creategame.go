package handlers

import (
	"github.com/BenStokmans/reversi-server/game"
	"github.com/sirupsen/logrus"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"
	"math/rand"
)

func handleCreateGame(msg *anypb.Any, client *game.Client) error {
	createGame := &game.CreateGame{}
	err := anypb.UnmarshalTo(msg, createGame, proto.UnmarshalOptions{})
	if err != nil {
		return err
	}

	resp := &game.CreateGameResponse{}
	if client.State.Game != nil {
		errMsg := "already in a game"
		resp.Error = &errMsg
		client.Send(resp)
		return nil
	}

	client.State.Game, err = game.NewGame(client)
	if err != nil {
		errMsg := err.Error()
		resp.Error = &errMsg
		client.Send(resp)
		return nil
	}
	client.State.Username = createGame.PlayerName

	if createGame.Color == game.Color_RANDOM {
		if rand.Intn(2) == 0 {
			resp.Color = game.Color_WHITE
		} else {
			resp.Color = game.Color_BLACK
		}
	} else {
		resp.Color = createGame.Color
	}

	err = client.State.Game.AddPlayer(client, resp.Color)
	if err != nil {
		errMsg := err.Error()
		resp.Error = &errMsg
		client.Send(resp)
		return nil
	}
	resp.GameId = int64(client.State.Game.Id)

	client.Send(resp)
	logrus.Debugf("Created game %d", client.State.Game.Id)
	return nil
}
