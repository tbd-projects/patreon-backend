package repository_postgresql

import "image/color"

const (
	AlphaByte = iota + 1
	BlueByte
	GreenByte
	RedByte
)

func getByte(number uint64, nByte uint64) uint8 {
	return uint8((number >> ((nByte - 1) * 8)) & ((1 << 9) - 1))
}

func setByte(number uint8, nByte uint64) uint64 {
	return uint64(number) << ((nByte - 1) * 8)
}

func convertRGBAToUint64(clr color.RGBA) uint64 {
	return setByte(clr.R, RedByte) +
		setByte(clr.B, BlueByte) +
		setByte(clr.G, GreenByte) +
		setByte(clr.A, AlphaByte)
}

func convertUint64ToRGBA(clr uint64) color.RGBA {
	return color.RGBA{R: getByte(clr, RedByte),
		B: getByte(clr, BlueByte),
		G: getByte(clr, GreenByte),
		A: getByte(clr, AlphaByte)}
}
