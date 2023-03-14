package game

type Game struct {
	id      uint64
	players []*Client
	state   *Board
}
