package main

import "math"

type Point struct {
	row, col uint8
}

func (s *Point) Left() bool {
	return s.col < 6
}

func (s *Point) Right() bool {
	return s.col >= 6
}

func (p *Point) HomeRow() bool {
	return p.row == 2
}

func (p *Point) Finger() uint8 {
	switch {
	case p.col < 5:
		return p.col
	case p.col == 5:
		return 4
	case p.col == 6:
		return 7
	case p.col <= 10:
		return p.col
	default:
		return 10
	}
}

func Dist2(p1, p2 Point) float64 {
	x := float64(p1.col - p2.col)
	y := float64(p1.row - p2.row)
	return x*x + y*y
}

func Dist(p1, p2 Point) float64 {
	return math.Sqrt(Dist2(p1, p2))
}

func PairEq(a, b, c, d uint8) bool {
	return a == c && b == d || a == d && b == c
}
