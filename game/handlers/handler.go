package handlers

import "google.golang.org/protobuf/types/known/anypb"

func HandleMessage(msg *anypb.Any) error {
	switch msg.MessageName() {
	case "reversi.CreateGame":
		return HandleCreateGame(msg)
	default:
		return nil
	}
}
