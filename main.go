package main

import "fmt"


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
	Eval(input string) float64
}

type Specimen struct {
	// permutation of some sort?
}

type Key struct {
	c rune
}

type Layout40 struct {
	
}

func main() {
	fmt.Println(L40Charset)
}
