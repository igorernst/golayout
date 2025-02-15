package main

import (
	"fmt"
	"log"
	"math"
	"math/rand"
	"sort"
	"sync"
)

const (
	PopulationSize uint = 50
	BestPartSize        = PopulationSize / 5
	HallOfFameSize uint = PopulationSize / 5
)

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

func dist2(p1, p2 Point) float64 {
	x := float64(p1.col - p2.col)
	y := float64(p1.row - p2.row)
	return x*x + y*y
}

func dist(p1, p2 Point) float64 {
	return math.Sqrt(dist2(p1, p2))
}

type genome struct {
	mapping map[rune]Point
	charset *[]rune
}

func (s *genome) PrettyPrint() {
	var (
		ar = new([12][3]rune)
	)
	for k, p := range s.mapping {
		ar[p.col-1][p.row-1] = k
	}
	for row := 1; row <= 3; row++ {
		for col := 1; col <= 13-row; col++ {
			fmt.Printf("%c ", ar[col-1][row-1])
		}
		fmt.Println()
	}
}

func (s *genome) Score(input string) float64 {
	const (
		homeRowBonus           = 0.5
		scissorsPenalty        = 2.0
		rowRedirectPenalty     = 2.0
		colRedirectPenalty     = 0.5
		sameColumnPenalty      = 1.0
		pinkyOffHomeRowPenalty = 0.5
	)
	var (
		prev   Point
		colInc bool
		rowInc bool
		score  float64 = 0
	)
	for _, r := range input {
		p, b := s.mapping[r]
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
		scissors := sameHand && pairEq(p.row, prev.row, 1, 3)
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

func pairEq(a, b, c, d uint8) bool {
	return a == c && b == d || a == d && b == c
}

func (s *genome) Crossover(s2 *genome) genome {
	var (
		usedPoints = make(map[Point]bool)
		sp         = make(map[rune]Point)
	)
	toFill := 0
	// take left from s
	for k, v := range s.mapping {
		toFill++
		if v.Left() {
			sp[k] = v
			usedPoints[v] = true
			toFill--
		}
	}
	// take right from s2, not taking any duplicates
	for k, v := range s2.mapping {
		if v.Left() {
			continue
		}
		if _, b := sp[k]; !b {
			sp[k] = v
			usedPoints[v] = true
			toFill--
		}
	}
	// finding letters/keys that are not assigned, filling them randomly
	runesNotUsed := make([]rune, toFill)
	pointsNotUsed := make([]Point, toFill)
	i := 0
	for k := range s.mapping {
		if _, b := sp[k]; !b {
			runesNotUsed[i] = k
			i++
		}
	}
	j := 0
	for _, v := range s.mapping {
		if _, b := usedPoints[v]; !b {
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
	return genome{
		mapping: sp,
		charset: s.charset,
	}
}

type instance struct {
	genome
	score float64
}

type Generation struct {
	population []instance
	hallOfFame []instance
}

func SeedGeneration() Generation {
	// TODO: fill this
	return Generation{
		population: []instance{},
		hallOfFame: []instance{},
	}
}

// TODO: show top 10 for each epoch (animated)

func TakeBest(s []instance) []instance {
	n := len(s)
	k := n / 5 // 20%
	sort.Slice(s, func(i, j int) bool {
		return s[i].score < s[j].score
	})
	if len(s[0:k]) != k {
		log.Panicln("wrong length of best", k)
	}
	return s[0:k]
}

func Extend(s, h []instance) error {
	k := int(BestPartSize + HallOfFameSize)
	filled := min((len(s)), k)
	for i := int(BestPartSize); i < filled; i++ {
		s[i] = h[i-int(HallOfFameSize)]
	}
	for i := filled; i < len(s); i++ {
		mutateOrCrossover := rand.Intn(2) // TODO: make a parameter/constant
		if mutateOrCrossover == 0 {
			a := rand.Intn(k)
			b := rand.Intn(k)
			s[i] = instance{
				genome: s[a].Crossover(&s[b].genome),
			}
		} else {
			a := rand.Intn(k)
			s[i] = s[a]
			s[i].Mutate1()
		}
	}
	return nil
	// TODO: should we take hof in here?
}

func (s *Generation) UpdateScores(input string) {
	var (
		wg sync.WaitGroup
	)
	n := len(s.population)
	wg.Add(n)
	for i := 0; i < n; i++ {
		go func(k int) {
			inst := s.population[k]
			inst.score = inst.Score(input)
			s.population[k] = inst
			wg.Done()
		}(i)
	}
	wg.Wait()
}

func OneStep(generation Generation, input string) Generation {
	generation.UpdateScores(input)
	tops := TakeBest(generation.population)                     // basically sort by score and ignore everything > top
	err := Extend(generation.population, generation.hallOfFame) // crossover and mutate
	if err != nil {
		log.Print("Error: extend, ", tops, generation.hallOfFame)
	}
	// TODO: implement TakeHOF / inplace variant
	hallOfFame := generation.hallOfFame // TakeHOF(generation)
	// TODO: make the function inplace
	return Generation{
		population: generation.population,
		hallOfFame: hallOfFame,
	}
}

type Layout40 struct {
}

func (s *genome) Mutate1() {
	// inplace elementary mutation
	charset := *s.charset
	n := len(charset)
	i := rand.Intn(n)
	j := rand.Intn(n)
	s.mapping[charset[i]], s.mapping[charset[j]] = s.mapping[charset[j]], s.mapping[charset[i]]
}

func main() {
}
