package main

import (
	"fmt"
	"log"
	"math/rand"
	"sort"
)

type genome struct {
	mapping map[rune]Point
	charset *[]rune
	hash    uint64
}

func (s *genome) Hash() {
	var (
		hash uint64 = 0
		m    uint64 = 1
	)
	keys := make([]rune, 0, len(s.mapping))
	for k := range s.mapping {
		keys = append(keys, k)
	}
	sort.Slice(keys, func(i int, j int) bool {
		return keys[i] > keys[j]
	})
	for _, k := range keys {
		v := s.mapping[k]
		hash += uint64(v.col+3) + uint64(v.row+6)*27
		hash *= m
		m *= 3
	}
	s.hash = hash
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
		homeRowBonus           = 0.065
		scissorsPenalty        = 2.0
		rowRedirectPenalty     = 2.0
		colRedirectPenalty     = 0.5
		sameFingerPenalty      = 2.0
		pinkyOffHomeRowPenalty = 0.5
	)
	var (
		prev        Point
		colInc      bool
		rowInc      bool
		score       float64 = 0
		fingerCount [11]int
	)
	for _, r := range input {
		p, b := s.mapping[r]
		if !b {
			continue
		}
		var penalty float64 = 0.0
		newColInc := p.col > prev.col
		newRowInc := p.row > prev.row
		//sameRow := prev.row == p.row
		sameFinger := prev.Finger() == p.Finger()
		sameHand := prev.Left() && p.Left() || prev.Right() && p.Right()
		colRedirect := sameHand && newColInc != colInc
		rowRedirect := sameHand && newRowInc != rowInc
		scissors := sameHand && PairEq(p.row, prev.row, 1, 3)
		pinkyOffHomeRow := !p.HomeRow() && (p.Finger() == 1 || p.Finger() == 10)
		if p.HomeRow() {
			score += homeRowBonus
		}
		if sameFinger && p != prev {
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
		fingerCount[p.Finger()]++
		rowInc = newRowInc
		colInc = newColInc
		prev = p
		score -= penalty
	}
	return score
}

func (s *genome) Crossover(s2 *genome, sp map[rune]Point) genome {
	var (
		usedPoints = make(map[Point]bool)
		usedKeys   = make(map[rune]bool)
	)
	toFill := 0
	// take left from s
	for k, v := range s.mapping {
		toFill++
		if v.Left() {
			sp[k] = v
			usedKeys[k] = true
			usedPoints[v] = true
			toFill--
		}
	}
	// take right from s2, not taking any duplicates
	for k, v := range s2.mapping {
		if v.Left() {
			continue
		}
		if _, b := usedKeys[k]; !b {
			sp[k] = v
			usedPoints[v] = true
			usedKeys[k] = true
			toFill--
		}
	}
	// finding letters/keys that are not assigned, filling them randomly
	runesNotUsed := make([]rune, toFill)
	pointsNotUsed := make([]Point, toFill)
	i := 0
	for k := range s.mapping {
		if _, b := usedKeys[k]; !b {
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
