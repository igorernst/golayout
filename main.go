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
		'+': {}
	}
)

// key to position transform
// or to a vertex on a graph

func textFilter(text string) []runes {
	
}

type Key struct {
	c rune
}

type Layout40 struct {
	
}

func main() {
	fmt.Println(L40Charset)
}

