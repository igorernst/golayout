package main

import (
	"fmt"
	"log"
	"math"
	"math/rand"
	//"golayout/constants"
)

const (
	PopulationSize uint = 50
	HallOfFameSize uint = 10
)

var (
	L40Charset = map[rune]struct{}{
		// lowercase only, '<' is on the same key as ','; numbers don't move etc.
		'a': {},
		'b': {},
		'c': {},
		'd': {},
		'e': {},
		'f': {},
		'g': {},
		'h': {},
		'i': {},
		'j': {},
		'k': {},
		'l': {},
		'm': {},
		'n': {},
		'o': {},
		'p': {},
		'q': {},
		'r': {},
		's': {},
		't': {},
		'u': {},
		'v': {},
		'w': {},
		'x': {},
		'y': {},
		'z': {},
		',': {},
		'<': {},
		'.': {},
		'>': {},
		'/': {},
		'?': {},
		'\'': {},
		'|': {},
		';': {},
		':': {},
		'[': {},
		'{': {},
		']': {},
		'}': {},
		'-': {},
		'_': {},
		'=': {},
		'+': {},
	}
)

// key to position transform
// or to a vertex on a graph

func textFilter(text string, charset map[rune]struct{}) []rune {
	// filters out the symbols that are not in the given charset
	// probably can work without filtering, cuz it created another array.
	// better just skip irrelevant chars, there are not a lot of those.
	// TODO actual filtering
	runes := []rune(text)
	return runes
}

type Specimen interface {
	Score(input string) float64
	Crossover(*Specimen) Specimen
}

type Point struct {
	row,col uint8
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

func dist2(p1,p2 Point) float64 {
	x := float64(p1.col - p2.col)
	y := float64(p1.row - p2.row)
	return x*x + y*y
}

func dist(p1,p2 Point) float64 {
	return math.Sqrt(dist2(p1, p2))
}

type specimen struct {
	mapping map[rune]Point
	// ignores graph structure. how to find bad pairs etc.
	// enumerate bad edges somehow, the graph is constant
	// same finger = same column
}

var (
	qwerty = map[rune]Point{
		'q': Point{1,1},
		'w': Point{1,2},
		'e': Point{1,3},
		'r': Point{1,4},
		't': Point{1,5},
		'y': Point{1,6},
		'u': Point{1,7},
		'i': Point{1,8},
		'o': Point{1,9},
		'p': Point{1,10},
		'[': Point{1,11},
		']': Point{1,12},
		'a': Point{2,1},
		's': Point{2,2},
		'd': Point{2,3},
		'f': Point{2,4},
		'g': Point{2,5},
		'h': Point{2,6},
		'j': Point{2,7},
		'k': Point{2,8},
		'l': Point{2,9},
		';': Point{2,10},
		'\'': Point{2,11},
		'z': Point{3,1},
		'x': Point{3,2},
		'c': Point{3,3},
		'v': Point{3,4},
		'b': Point{3,5},
		'n': Point{3,6},
		'm': Point{3,7},
		',': Point{3,8},
		'.': Point{3,9},
		'/': Point{3,10},
	}
	nerps = map[rune]Point{
		'x': Point{1,1},
		'l': Point{1,2},
		'd': Point{1,3},
		'p': Point{1,4},
		'v': Point{1,5},
		'z': Point{1,6},
		'k': Point{1,7},
		'o': Point{1,8},
		'u': Point{1,9},
		';': Point{1,10},
		'[': Point{1,11},
		']': Point{1,12},
		'n': Point{2,1},
		'r': Point{2,2},
		't': Point{2,3},
		's': Point{2,4},
		'g': Point{2,5},
		'y': Point{2,6},
		'h': Point{2,7},
		'e': Point{2,8},
		'i': Point{2,9},
		'a': Point{2,10},
		'/': Point{2,11},
		'j': Point{3,1},
		'm': Point{3,2},
		'c': Point{3,3},
		'w': Point{3,4},
		'q': Point{3,5},
		'b': Point{3,6},
		'f': Point{3,7},
		'\'': Point{3,8},
		',': Point{3,9},
		'.': Point{3,10},
	}
)

var blockedEdges map[int]bool = map[int]bool{}
// if else is (a,b) to (c,d) we write into
// a + b >> 8 + c >> 16 + d >> 24  which is an int

func (s *specimen) Score(input string) float64 {
	const (
		homeRowBonus = 0.5
		scissorsPenalty = 2.0
		rowRedirectPenalty = 2.0
		colRedirectPenalty = 0.5
		sameColumnPenalty = 1.0
		pinkyOffHomeRowPenalty = 0.5
	)
	var (
		prev Point
		colInc bool
		rowInc bool
		score float64 = 0
	)
	for _, r := range input {
		p,b := s.mapping[r];
		if !b {
			continue
		}
		var penalty float64 = 0.0
		newColInc := p.col > prev.col
		newRowInc := p.row > prev.row
		sameRow := prev.row == p.row
		sameColumn := prev.col == p.col
		sameHand := prev.Left() && p.Left() || prev.Right() && p.Right()
		colRedirect := sameHand && newColInc != colInc
		rowRedirect := sameHand && newRowInc != rowInc
		scissors := sameHand && pairEq(p.row,prev.row,1,3)
		if p.HomeRow() {
			score += homeRowBonus
		}
		if sameColumn && !sameRow {
			penalty += sameColumnPenalty
		}
		if colRedirect {
			penalty += colRedirectPenalty
		}
		if rowRedirect {
			penalty += rowRedirectPenalty
		}
		if scissors {
			penalty += scissorsPenalty
		}
		rowInc = newRowInc
		colInc = newColInc
		prev = p
		score -= penalty
	}
	return score
}

func pairEq(a,b,c,d uint8) bool {
	return a==c && b==d || a==d && b==c
}

func (s *specimen) Crossover(s2 *specimen) specimen {
	var (
		usedPoints = make(map[Point]bool)
		sp = make(map[rune]Point)
	)
	toFill := 0
	// take left from s
	for k,v := range s.mapping {
		toFill++
		if v.Left() {
			sp[k] = v
			usedPoints[v] = true
			toFill--
		}
	}
	// take right from s2, not taking any duplicates
	for k,v := range s2.mapping {
		if v.Left() {
			continue
		}
		if _,b := sp[k]; !b {
			sp[k] = v
			usedPoints[v] = true
			toFill--
		}
	}
	// finding letters/keys that are not assigned
	runesNotUsed := make([]rune,toFill)
	pointsNotUsed := make([]Point,toFill)
	i := 0
	for k := range s.mapping {
		if _,b := sp[k]; !b {
			runesNotUsed[i] = k
			i++
		}
	}
	j := 0
	for _,v := range s.mapping {
		if _,b := usedPoints[v]; !b {
			pointsNotUsed[j] = v
			j++
		}
	}
	if i != j {
		log.Panicln("i not equal to j", i, j)
	}
	if i != toFill {
		log.Panicln("i not equal to toFill", i, toFill)
	}
	rand.Shuffle(len(runesNotUsed),
		func(i, j int) { runesNotUsed[i], runesNotUsed[j] = runesNotUsed[j], runesNotUsed[i] })
	for k := 0; k < i; k++ {
		sp[runesNotUsed[k]] = pointsNotUsed[k]
	}
	// TODO: consider taking the right first if that makes more keys before random filling
	if len(sp) != len(s.mapping) {
		panic("The output length is not equal to the input")
	}
	return specimen{
		mapping: sp,
	}
}

type Generation struct {
	population []Specimen
	hallOfFame []Specimen
}

func SeedGeneration() Generation {
	// TODO: fill this
	return Generation{
		population: []Specimen{},
		hallOfFame: []Specimen{},
	}
}

/*
func OneStep(generation Generation) Generation {
	g := Extend(generation) // crossover and mutate
	population := TakeBest(g)
	hallOfFame := TakeHOF(generation.hallOfFame, g)
	return Generation{
		population: population,
		hallOfFame: hallOfFame,
	}
        }
*/

type Layout40 struct {

}

func crossover_test() {
	var (
		qwerty = specimen{
			mapping: qwerty,
		}
		nerps = specimen{
			mapping: nerps,
		}
		s = "hello from the testing string"
	)
	s1 := qwerty.Crossover(&nerps)
	fmt.Println(qwerty.mapping, "\n", nerps.mapping, "\n", s1.mapping)
	fmt.Println(qwerty.Score(s), nerps.Score(s), s1.Score(s))
}

func main() {
}
