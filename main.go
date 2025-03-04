package main

import (
	"fmt"
	"os"
)

func main() {
	bytes, err := os.ReadFile("text/Alice.txt")
	if err != nil {
		fmt.Println(err)
		return
	}
	input := string(bytes)
	gen := SeedGeneration()
	gen.UpdateScores(input[0:1000])
	l := len(input)
	fmt.Println("input length: ", l)
	for i := 0; i < 100; i++ {
		gen = OneStep(gen, input)
	}
	fmt.Println("Hall Of Fame:")
	for i, v := range gen.hallOfFame {
		fmt.Println(i, ": ", v.score, ", hash:", v.hash)
		v.genome.PrettyPrint()
	}
}
