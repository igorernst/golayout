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
}


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
	// finding letters/keys that are not assigned, filling them randomly
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
			mapping: Qwerty,
		}
		nerps = specimen{
			mapping: Nerps,
		}
		s = "hello from the testing string"
	)
	s1 := qwerty.Crossover(&nerps)
	fmt.Println(qwerty.mapping, "\n", nerps.mapping, "\n", s1.mapping)
	fmt.Println(qwerty.Score(s), nerps.Score(s), s1.Score(s))
}

func main() {
}
