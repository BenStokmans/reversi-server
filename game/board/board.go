package board

import "errors"

type Board struct {
	player   uint64
	opponent uint64
}

func NewBoard() *Board {
	return &Board{
		player:   0x0000001008000000,
		opponent: 0x0000000810000000,
	}
}

func (b *Board) Move(x, y uint8) error {
	place := y*8 + x
	moves := legalMoves(b.player, b.opponent)
	if moves&(1<<place) == 0 {
		return errors.New("illegal move")
	}
	flip := calcFlip(b.player, b.opponent, x, y)
	b.player |= flip
	b.opponent &= ^flip
	b.SwitchTurn()
	return nil
}

func (b *Board) LegalMoves() uint64 {
	return legalMoves(b.player, b.opponent)
}

func (b *Board) SwitchTurn() {
	b.player, b.opponent = b.opponent, b.player
}

func (b *Board) PlayerScore() uint8 {
	return countSetBits(b.player)
}

func (b *Board) OpponentScore() uint8 {
	return countSetBits(b.opponent)
}

func (b *Board) GameOver() bool {
	return legalMoves(b.player, b.opponent) == 0 && legalMoves(b.opponent, b.player) == 0
}
