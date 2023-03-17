package board

var flipPreCalc [64][8][7]uint64

func init() {
	for y := uint8(0); y < 8; y++ {
		for x := uint8(0); x < 8; x++ {
			initPreCalcPart(x, y)
		}
	}
}

func initPreCalcPart(x, y uint8) {
	place := y*8 + x
	di := uint8(0)

	line := uint64(0)
	li := uint8(0)
	for dy := int8(-1); dy <= 1; dy++ {
		for dx := int8(-1); dx <= 1; dx++ {
			if dx == 0 && dy == 0 {
				continue
			}
			line = 0
			li = 0
			cx := int8(x) + dx
			cy := int8(y) + dy

			for (cx >= 0 && cx < 8) && (cy >= 0 && cy < 8) {
				line |= uint64(1) << (cy*8 + cx)
				flipPreCalc[place][di][li] = line
				cx += dx
				cy += dy
				li++
			}
			di++
		}
	}
}

func getLinesFromPoint(x, y uint8) (out uint64) {
	place := y*8 + x
	for i := uint8(0); i < 8; i++ {
		for j := uint8(0); j < 7; j++ {
			if flipPreCalc[place][i][j] == 0 {
				break
			}
			out |= flipPreCalc[place][i][j]
		}
	}
	return
}

func calcFlip(p, o uint64, x, y uint8) (out uint64) {
	place := y*8 + x
	open := ^(p | o)
	for i := uint8(0); i < 8; i++ {
		of := false
		for j := uint8(0); j < 7; j++ {
			v := flipPreCalc[place][i][j]
			if v == 0 {
				break
			}
			if open&v != 0 {
				break
			}
			if (o&v) != 0 && !of {
				of = true
				continue
			}
			if (p&v) != 0 && !of {
				break
			}
			if (p&v) != 0 && of {
				out |= flipPreCalc[place][i][j-1]
				break
			}
		}
	}
	out |= uint64(1) << place
	return
}
