package main

import (
	"testing"
)

func crossover_test(t *testing.T) {
	var (
		qwerty = genome{
			mapping: Qwerty,
			charset: &StandardCharset,
		}
		nerps = genome{
			mapping: Nerps,
			charset: &StandardCharset,
		}
	)
	s1 := qwerty.Crossover(&nerps)
	if len(s1.mapping) != len(qwerty.mapping) {
		t.Fail()
	}
}

func mutate1_test(t *testing.T) {
	var (
		qwerty = genome{
			mapping: Qwerty,
			charset: &StandardCharset,
		}
	)
	qwerty.Mutate1()
	t.SkipNow()
}

func score_test(t *testing.T) {
	var (
		qwerty = genome{
			mapping: Qwerty,
			charset: &StandardCharset,
		}
		s = "hello from the test string"
	)
	qwerty.Score(s)
}
