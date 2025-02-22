package main

import (
	"fmt"
	"log"
	"math"
	"math/rand"
	"os"
	"sort"
	"sync"
)

const (
	PopulationSize    uint = 50
	BestPartSize           = PopulationSize / 5
	HallOfFameSize    uint = PopulationSize / 5
	MutationFrequency      = 10 // mutataion would happen n-1 time out of n
	// TODO: make the constants configurable using a struct
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
			fmt.Printf("%c", ar[col-1][row-1])
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
		sameFingerPenalty      = 2.0
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
		sameFinger := prev.Finger() == p.Finger()
		sameHand := prev.Left() && p.Left() || prev.Right() && p.Right()
		colRedirect := sameHand && newColInc != colInc
		rowRedirect := sameHand && newRowInc != rowInc
		scissors := sameHand && pairEq(p.row, prev.row, 1, 3)
		pinkyOffHomeRow := !p.HomeRow() && (p.Finger() == 1 || p.Finger() == 10)
		if p.HomeRow() {
			score += homeRowBonus
		}
		if sameFinger && !sameRow {
			penalty += sameFingerPenalty
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
		if pinkyOffHomeRow {
			penalty += pinkyOffHomeRowPenalty
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
	s := make([]instance, PopulationSize)
	for i := 0; i < int(PopulationSize); i++ {
		s[i].mapping = make(map[rune]Point)
		s[i].charset = &StandardCharset
	}
	h := make([]instance, HallOfFameSize)
	for i := 0; i < int(HallOfFameSize); i++ {
		h[i].mapping = make(map[rune]Point)
		h[i].charset = &StandardCharset
	}
	for k, v := range Qwerty.mapping {
		s[0].mapping[k] = v
	}
	for k, v := range Nerps.mapping {
		s[1].mapping[k] = v
	}
	m := 2
	for i := 0; i < int(HallOfFameSize); i++ {
		for k, v := range s[min(i, m-1)].mapping {
			h[i].mapping[k] = v
		}
	}
	k := int(m)
	filled := min(len(s), k)
	for i := filled; i < len(s); i++ {
		mutateOrCrossover := rand.Intn(MutationFrequency)
		if mutateOrCrossover == 0 {
			a := rand.Intn(k)
			b := rand.Intn(k)
			s[i] = instance{
				genome: s[a].Crossover(&s[b].genome), // TODO: change to overwriting rather than making a new entity
			}
		} else {
			a := rand.Intn(k)
			if s[a].charset == nil {
				fmt.Println("wrong element", a, s[a])
			}
			g := genome{
				mapping: make(map[rune]Point), // TODO: same
				charset: s[a].charset,
			}
			for k, v := range s[a].mapping {
				g.mapping[k] = v
			}
			g.Mutate1()
			s[i] = instance{genome: g}
		}
	}
	g := Generation{
		population: s,
		hallOfFame: h,
	}
	return g
}

// TODO: show top 10 for each epoch (animated)

func TakeBest(s []instance) []instance {
	n := len(s)
	k := n / 5 // 20%
	sort.Slice(s, func(i, j int) bool {
		return s[i].score > s[j].score
	})
	if len(s[0:k]) != k {
		log.Panicln("wrong length of best", k)
	}
	return s[0:k]
}

func Extend(s, h []instance) error {
	k := int(BestPartSize + HallOfFameSize)
	filled := min(len(s), k)
	for i := int(BestPartSize); i < filled; i++ {
		for k, v := range h[i-int(HallOfFameSize)].mapping {
			s[i].mapping[k] = v
		}
		s[i].charset = h[i-int(HallOfFameSize)].charset
	}
	for i := filled; i < len(s); i++ {
		mutateOrCrossover := rand.Intn(MutationFrequency)
		if mutateOrCrossover == 0 {
			a := rand.Intn(k)
			b := rand.Intn(k)
			s[i] = instance{
				genome: s[a].Crossover(&s[b].genome), // TODO: change to overwriting rather than making a new entity
			}
		} else {
			a := rand.Intn(k)
			if s[a].charset == nil {
				fmt.Println("Extend: wrong element", a, s[a].genome)
			}
			g := genome{
				mapping: make(map[rune]Point), // TODO: same
				charset: s[a].charset,
			}
			for k, v := range s[a].mapping {
				g.mapping[k] = v
			}
			g.Mutate1()
			s[i] = instance{genome: g}
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
	m := len(s.hallOfFame)
	wg.Add(n + m)
	for i := 0; i < n; i++ {
		go func(k int) {
			s.population[k].score = s.population[k].Score(input)
			wg.Done()
		}(i)
	}
	for i := 0; i < m; i++ {
		go func(k int) {
			s.hallOfFame[k].score = s.hallOfFame[k].Score(input)
			wg.Done()
		}(i)
	}
	wg.Wait()
}

func OneStep(generation Generation, input string) Generation {
	tops := TakeBest(generation.population)                     // basically sort by score and ignore everything > top
	err := Extend(generation.population, generation.hallOfFame) // crossover and mutate
	if err != nil {
		log.Print("Error: extend, ", tops, generation.hallOfFame)
	}
	generation.UpdateScores(input)
	// update hallOfFame with the top of the population
	sort.Slice(generation.hallOfFame, func(i, j int) bool {
		return generation.hallOfFame[i].score > generation.hallOfFame[j].score
	})
	l := 0
	r := 0
	comp := func(l, r int) bool {
		return generation.hallOfFame[r].score > generation.population[l].score
	}
	for l+r < int(HallOfFameSize) {
		if comp(l, r) {
			r++
			continue
		}
		// save r-th to the bottom
		for k, v := range generation.hallOfFame[r].mapping {
			generation.hallOfFame[int(HallOfFameSize)-l-1].mapping[k] = v
		}
		// move population[l] to r-th place
		for k, v := range generation.population[l].mapping {
			generation.hallOfFame[r].mapping[k] = v
		}
		l++
	}
	sort.Slice(generation.hallOfFame, func(i, j int) bool {
		return generation.hallOfFame[i].score > generation.hallOfFame[j].score
	})
	return generation
}

type Layout40 struct {
}

func (s *genome) Mutate1() {
	// inplace elementary mutation
	if s.charset == nil {
		s.PrettyPrint()
	}
	charset := *s.charset
	n := len(charset)
	i := rand.Intn(n)
	j := rand.Intn(n)
	s.mapping[charset[i]], s.mapping[charset[j]] = s.mapping[charset[j]], s.mapping[charset[i]]
}

func main() {
	// TODO: read file
	bytes, err := os.ReadFile("text/Alice.txt")
	if err != nil {
		fmt.Println(err)
		return
	}
	input := string(bytes)
	seed := SeedGeneration()
	seed.UpdateScores(input[0:1000])
	/*for i, v := range seed.population {
		fmt.Println(i)
		v.genome.PrettyPrint()
	}*/
	gen := seed
	l := len(input)
	fmt.Println("input length: ", l)
	for i := 0; i < 200; i++ {
		a := i * 1000
		b := a + 10000
		if b < l {
			gen = OneStep(gen, input[a:b])
		}
	}
	for i, v := range gen.population {
		fmt.Println(i, ": ", v.score)
		v.genome.PrettyPrint()
	}
	fmt.Println("Hall Of Fame:")
	for i, v := range gen.hallOfFame {
		fmt.Println(i, ": ", v.score)
		v.genome.PrettyPrint()
	}
	// TODO: Hall of fame is not updated properly
	// TODO: Make the scope update in the end? Makes sense
}
