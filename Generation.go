package main

import (
	"fmt"
	"log"
	"math/rand"
	"sort"
	"sync"
)

const (
	PopulationSize    uint = 100
	BestPartSize           = PopulationSize / 10
	HallOfFameSize    uint = PopulationSize / 10
	MutationFrequency      = 5 // mutataion would happen n-1 time out of n
	// TODO: make the constants configurable using a struct
)

type instance struct {
	genome
	score float64
}

type Generation struct {
	population []instance
	hallOfFame []instance
}

func SeedGeneration() Generation {
	s := make([]instance, PopulationSize)
	hashes := make(map[uint64]struct{}, len(s))
	for i := 0; i < int(PopulationSize); i++ {
		s[i].mapping = make(map[rune]Key, 35)
		s[i].charset = &StandardCharset
	}
	h := make([]instance, HallOfFameSize)
	for i := 0; i < int(HallOfFameSize); i++ {
		h[i].mapping = make(map[rune]Key, 35)
		h[i].charset = &StandardCharset
	}
	for k, v := range Qwerty.mapping {
		s[0].mapping[k] = v
	}
	for k, v := range Nerps.mapping {
		s[1].mapping[k] = v
	}
	m := 2
	for i := 0; i < m; i++ {
		s[i].Hash()
		hashes[s[i].hash] = struct{}{}
	}
	for i := 0; i < int(HallOfFameSize); i++ {
		for k, v := range s[min(i, m-1)].mapping {
			h[i].mapping[k] = v
		}
		h[i].Hash()
	}
	k := int(m)
	filled := min(len(s), k)
	for i := filled; i < len(s); i++ {
		for {
			mutateOrCrossover := rand.Intn(MutationFrequency)
			if mutateOrCrossover == 0 {
				a := rand.Intn(k)
				b := rand.Intn(k)
				s[a].Crossover(&s[b].genome, s[i].mapping)
				s[i].Hash()
				if _, b := hashes[s[i].hash]; !b {
					hashes[s[i].hash] = struct{}{}
					break
				}

			} else {
				a := rand.Intn(k)
				if s[a].charset == nil {
					fmt.Println("wrong element", a, s[a])
				}
				g := genome{
					mapping: make(map[rune]Key, 35),
					charset: s[a].charset,
				}
				for k, v := range s[a].mapping {
					g.mapping[k] = v
				}
				g.Mutate1()
				s[i] = instance{genome: g}
				s[i].Hash()
				if _, b := hashes[s[i].hash]; !b {
					hashes[s[i].hash] = struct{}{}
					break
				}
			}
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
	hashes := make(map[uint64]struct{}, len(s))
	filled := min(len(s), k)
	for i := int(BestPartSize); i < filled; i++ {
		for k, v := range h[i-int(HallOfFameSize)].mapping {
			s[i].mapping[k] = v
		}
		s[i].charset = h[i-int(HallOfFameSize)].charset
		s[i].Hash()
		hashes[s[i].hash] = struct{}{}
	}
	for i := filled; i < len(s); i++ {
		for {
			mutateOrCrossover := rand.Intn(MutationFrequency)
			if mutateOrCrossover == 0 {

				a := rand.Intn(k)
				b := rand.Intn(k)
				s[i] = instance{
					genome: s[a].Crossover(&s[b].genome, s[i].mapping),
				}
				s[i].Hash()
				// if hash is used try once more
				if _, b := hashes[s[i].hash]; !b {
					hashes[s[i].hash] = struct{}{}
					break
				}

			} else {
				a := rand.Intn(k)
				if s[a].charset == nil {
					fmt.Println("Extend: wrong element", a, s[a].genome)
				}
				for k, v := range s[a].mapping {
					s[i].mapping[k] = v
				}
				s[i].genome.Mutate1()
				s[i].Hash()
				if _, b := hashes[s[i].hash]; !b {
					hashes[s[i].hash] = struct{}{}
					break
				}
			}
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
			s.population[k].Hash()
			wg.Done()
		}(i)
	}
	for i := 0; i < m; i++ {
		go func(k int) {
			s.hallOfFame[k].score = s.hallOfFame[k].Score(input)
			s.hallOfFame[k].Hash()
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
	hashes := make(map[uint64]struct{}, int(HallOfFameSize))
	sort.Slice(generation.hallOfFame, func(i, j int) bool {
		return generation.hallOfFame[i].score > generation.hallOfFame[j].score
	})
	for i := 0; i < int(HallOfFameSize); i++ {
		hashes[generation.hallOfFame[i].hash] = struct{}{}
	}
	l := 0
	r := 0
	comp := func(l, r int) bool {
		return generation.hallOfFame[r].score > generation.population[l].score
	}
	c := 0
	for c < int(HallOfFameSize) && l < len(generation.population) {
		if comp(l, r) {
			r++
			c++
			continue
		}
		if _, b := hashes[generation.population[l].hash]; b {
			l++
			// check next from the left side
			continue
		}
		hashes[generation.population[l].hash] = struct{}{}
		// save r-th to the bottom
		generation.hallOfFame[int(HallOfFameSize)-c-1].score = generation.hallOfFame[r].score
		for k, v := range generation.hallOfFame[r].mapping {
			generation.hallOfFame[int(HallOfFameSize)-c-1].mapping[k] = v
		}
		generation.hallOfFame[int(HallOfFameSize)-c-1].Hash()
		// move population[l] to r-th place
		generation.hallOfFame[r].score = generation.population[l].score
		for k, v := range generation.population[l].mapping {
			generation.hallOfFame[r].mapping[k] = v
		}
		generation.hallOfFame[r].Hash()
		l++
		c++
	}
	sort.Slice(generation.hallOfFame, func(i, j int) bool {
		return generation.hallOfFame[i].score > generation.hallOfFame[j].score
	})
	return generation
}
