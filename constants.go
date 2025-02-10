package main

var (
	Qwerty = map[rune]Point{
		'q': Point{1,1},
		'w': Point{1,2},
		'e': Point{1,3},
		'r': Point{1,4},
		't': Point{1,5},
		'y': Point{1,6},
		'u': Point{1,7},
		'i': Point{1,8},
		'o': Point{1,9},
		'p': Point{1,10},
		'[': Point{1,11},
		']': Point{1,12},
		'a': Point{2,1},
		's': Point{2,2},
		'd': Point{2,3},
		'f': Point{2,4},
		'g': Point{2,5},
		'h': Point{2,6},
		'j': Point{2,7},
		'k': Point{2,8},
		'l': Point{2,9},
		';': Point{2,10},
		'\'': Point{2,11},
		'z': Point{3,1},
		'x': Point{3,2},
		'c': Point{3,3},
		'v': Point{3,4},
		'b': Point{3,5},
		'n': Point{3,6},
		'm': Point{3,7},
		',': Point{3,8},
		'.': Point{3,9},
		'/': Point{3,10},
	}
	Nerps = map[rune]Point{
		'x': Point{1,1},
		'l': Point{1,2},
		'd': Point{1,3},
		'p': Point{1,4},
		'v': Point{1,5},
		'z': Point{1,6},
		'k': Point{1,7},
		'o': Point{1,8},
		'u': Point{1,9},
		';': Point{1,10},
		'[': Point{1,11},
		']': Point{1,12},
		'n': Point{2,1},
		'r': Point{2,2},
		't': Point{2,3},
		's': Point{2,4},
		'g': Point{2,5},
		'y': Point{2,6},
		'h': Point{2,7},
		'e': Point{2,8},
		'i': Point{2,9},
		'a': Point{2,10},
		'/': Point{2,11},
		'j': Point{3,1},
		'm': Point{3,2},
		'c': Point{3,3},
		'w': Point{3,4},
		'q': Point{3,5},
		'b': Point{3,6},
		'f': Point{3,7},
		'\'': Point{3,8},
		',': Point{3,9},
		'.': Point{3,10},
	}

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
