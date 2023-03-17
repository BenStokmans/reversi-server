package board

func countSetBits(n uint64) uint8 {
	n = n - ((n >> 1) & 0x5555555555555555)
	n = (n & 0x3333333333333333) + ((n >> 2) & 0x3333333333333333)
	return (uint8)(((n + (n >> 4)) & 0xF0F0F0F0F0F0F0F) * 0x101010101010101 >> 56)
}
