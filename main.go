package main

import (
	"fmt"
	"os"
)

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
	// for i := 0; i < 2000; i++ {
	// 	a := (i % 300) * 500
	// 	b := a + 5000
	// 	if b < l {
	// 		gen = OneStep(gen, input[a:b])
	// 	} else {
	// 		fmt.Printf("Generations passed: %d\n", i)
	// 		break
	// 	}
	// }
	for i := 0; i < 100; i++ {
		gen = OneStep(gen, input)
	}

	//for i, v := range gen.population {
	//	fmt.Println(i, ": ", v.score, ", hash:", v.hash)
	//	v.genome.PrettyPrint()
	//}
	fmt.Println("Hall Of Fame:")
	for i, v := range gen.hallOfFame {
		fmt.Println(i, ": ", v.score, ", hash:", v.hash)
		v.genome.PrettyPrint()
	}
	// TODO: Make the scope update in the end? Makes sense
	// TODO: improve scoring model
}
