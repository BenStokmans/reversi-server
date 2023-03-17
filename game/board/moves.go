package board

func legalMoves(player, opponent uint64) (moves uint64) {
	var flip1, flip7, flip9, flip8, pre1, pre7, pre9, pre8 uint64
	mO := opponent & 0x7e7e7e7e7e7e7e7e
	flip1 = mO & (player << 1)
	flip7 = mO & (player << 7)
	flip9 = mO & (player << 9)
	flip8 = opponent & (player << 8)
	flip1 |= mO & (flip1 << 1)
	flip7 |= mO & (flip7 << 7)
	flip9 |= mO & (flip9 << 9)
	flip8 |= opponent & (flip8 << 8)
	pre1 = mO & (mO << 1)
	pre7 = mO & (mO << 7)
	pre9 = mO & (mO << 9)
	pre8 = opponent & (opponent << 8)
	flip1 |= pre1 & (flip1 << 2)
	flip7 |= pre7 & (flip7 << 14)
	flip9 |= pre9 & (flip9 << 18)
	flip8 |= pre8 & (flip8 << 16)
	flip1 |= pre1 & (flip1 << 2)
	flip7 |= pre7 & (flip7 << 14)
	flip9 |= pre9 & (flip9 << 18)
	flip8 |= pre8 & (flip8 << 16)
	moves = flip1 << 1
	moves |= flip7 << 7
	moves |= flip9 << 9
	moves |= flip8 << 8
	flip1 = mO & (player >> 1)
	flip7 = mO & (player >> 7)
	flip9 = mO & (player >> 9)
	flip8 = opponent & (player >> 8)
	flip1 |= mO & (flip1 >> 1)
	flip7 |= mO & (flip7 >> 7)
	flip9 |= mO & (flip9 >> 9)
	flip8 |= opponent & (flip8 >> 8)
	pre1 >>= 1
	pre7 >>= 7
	pre9 >>= 9
	pre8 >>= 8
	flip1 |= pre1 & (flip1 >> 2)
	flip7 |= pre7 & (flip7 >> 14)
	flip9 |= pre9 & (flip9 >> 18)
	flip8 |= pre8 & (flip8 >> 16)
	flip1 |= pre1 & (flip1 >> 2)
	flip7 |= pre7 & (flip7 >> 14)
	flip9 |= pre9 & (flip9 >> 18)
	flip8 |= pre8 & (flip8 >> 16)
	moves |= flip1 >> 1
	moves |= flip7 >> 7
	moves |= flip9 >> 9
	moves |= flip8 >> 8
	moves &= ^(player | opponent)
	return
}
