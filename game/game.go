package game

import (
	"errors"
	"github.com/BenStokmans/reversi-server/game/board"
	"github.com/BenStokmans/reversi-server/snowflake"
)

type Game struct {
	owner *Server

	Id      snowflake.Snowflake
	Started bool
	Turn    Color
	white   *Client
	black   *Client
	state   *board.Board
}

func NewGame(c *Client) (*Game, error) {
	c.owner.gamesMut.Lock()
	defer c.owner.gamesMut.Unlock()
	for _, game := range c.owner.games {
		if game.white == c || game.black == c {
			return nil, errors.New("already in a game")
		}
	}

	g := &Game{
		owner: c.owner,
		Id:    snowflake.Next(),
		state: board.NewBoard(),
		Turn:  Color_BLACK,
	}
	g.owner.games[g.Id] = g

	return g, nil
}

func ClientJoin(c *Client, id snowflake.Snowflake) (Color, error) {
	c.owner.gamesMut.Lock()
	defer c.owner.gamesMut.Unlock()
	for _, game := range c.owner.games {
		if game.white == c || game.black == c {
			return 0, errors.New("already in a game")
		}
	}
	game, ok := c.owner.games[id]
	if !ok {
		return 0, errors.New("game does not exist")
	}

	color := Color_BLACK
	if game.white == nil {
		color = Color_WHITE
	}

	if err := game.AddPlayer(c, color); err != nil {
		return 0, err
	}
	c.State.Game = game
	return color, nil
}

func (g *Game) AddPlayer(c *Client, color Color) error {
	if g.white != nil && g.black != nil {
		return errors.New("game is full")
	}
	if c == g.white || c == g.black {
		return errors.New("already in game")
	}

	if g.white == nil && color == Color_WHITE {
		g.white = c
	} else if g.black == nil && color == Color_BLACK {
		g.black = c
	} else {
		return errors.New("color already taken")
	}
	return nil
}

func (g *Game) RemovePlayer(c *Client) {
	if g.white == c {
		g.white = nil
	}
	if g.black == c {
		g.black = nil
	}
	g.End()
}

func (g *Game) Move(c *Client, move *PlayMove) error {
	if !g.Started {
		return errors.New("game is not started")
	}
	if g.white != c && g.black != c {
		return errors.New("not in game")
	}
	if g.Turn == Color_WHITE && g.white != c {
		return errors.New("not client's turn")
	}
	if g.Turn == Color_BLACK && g.black != c {
		return errors.New("not client's turn")
	}

	if err := g.state.Move(uint8(move.X), uint8(move.Y)); err != nil {
		return err
	}
	gameMove := &GameMove{
		X: move.X,
		Y: move.Y,
	}
	// check if the current player has to pass
	if !g.state.GameOver() && g.state.LegalMoves() == 0 {
		g.state.SwitchTurn()
		gameMove.Pass = true
	}

	if g.white != c {
		g.white.Send(gameMove)
	}
	if g.black != c {
		g.black.Send(gameMove)
	}
	if g.state.GameOver() {
		g.End()
		return nil
	}

	// if we pass we don't need to switch turns
	if gameMove.Pass {
		return nil
	}

	if g.Turn == Color_WHITE {
		g.Turn = Color_BLACK
	} else {
		g.Turn = Color_WHITE
	}

	return nil
}

func (g *Game) Start() error {
	if g.white == nil || g.black == nil {
		return errors.New("game is not full")
	}
	g.white.Send(&GameStart{})
	g.black.Send(&GameStart{})
	g.Started = true
	return nil
}

func (g *Game) End() {
	if !g.Started {
		return
	}

	gameEnd := g.getGameEnded()
	if g.white != nil {
		g.white.Send(gameEnd)
		g.white.State.Game = nil
	}
	if g.black != nil {
		g.black.Send(gameEnd)
		g.black.State.Game = nil
	}

	g.Started = false
	g.owner.gamesMut.Lock()
	delete(g.owner.games, g.Id)
	g.owner.gamesMut.Unlock()
}

func (g *Game) getGameEnded() *GameEnded {
	gameEnd := &GameEnded{}
	var winner Color
	var reason string
	if g.white == nil {
		winner = Color_BLACK
		reason = "white player left"
		gameEnd.Winner = &winner
		gameEnd.SpecialReason = &reason
		return gameEnd
	}
	if g.black == nil {
		winner = Color_WHITE
		reason = "black player left"
		gameEnd.Winner = &winner
		gameEnd.SpecialReason = &reason
		return gameEnd
	}

	if !g.state.GameOver() {
		reason = "unknown"
		gameEnd.SpecialReason = &reason
		return gameEnd
	}
	if g.state.PlayerScore() > g.state.OpponentScore() {
		winner = Color_BLACK
		gameEnd.Winner = &winner
	} else if g.state.PlayerScore() < g.state.OpponentScore() {
		winner = Color_WHITE
		gameEnd.Winner = &winner
	} else {
		reason = "draw"
	}
	gameEnd.SpecialReason = &reason

	if gameEnd.Winner == nil {
		return gameEnd
	}

	if g.Turn == Color_BLACK {
		winner = Color_WHITE
		gameEnd.Winner = &winner
	} else {
		winner = Color_BLACK
		gameEnd.Winner = &winner
	}
	return gameEnd
}
