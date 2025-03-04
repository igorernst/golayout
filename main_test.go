package main

import (
	"testing"
)

func TestCrossover(t *testing.T) {
	g := genome{
		mapping: make(map[rune]Key),
		charset: Qwerty.charset,
	}
	s1 := Qwerty.Crossover(&Nerps, g.mapping)
	if len(s1.mapping) != len(Qwerty.mapping) {
		t.Fail()
	}
	s1.Mutate1()
}

func TestMutate1(t *testing.T) {
	Qwerty.Mutate1()
}

func TestScore(t *testing.T) {
	var (
		s = "hello from the test string"
	)
	Qwerty.Score(s)
}
