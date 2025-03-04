package main

var (
	Qwerty = genome{
		mapping: qwerty,
		charset: &StandardCharset,
	}
	Nerps = genome{
		mapping: nerps,
		charset: &StandardCharset,
	}
)

var (
	StandardCharset = []rune{'q', 'w', 'e', 'r', 't', 'y', 'u', 'i', 'o', 'p', '[', ']', 'a', 's', 'd', 'f', 'g', 'h', 'j', 'k', 'l', ';', '\'', 'z', 'x', 'c', 'v', 'b', 'n', 'm', ',', '.', '/'}
	qwerty          = map[rune]Key{
		'q':  Key{1, 1},
		'w':  Key{1, 2},
		'e':  Key{1, 3},
		'r':  Key{1, 4},
		't':  Key{1, 5},
		'y':  Key{1, 6},
		'u':  Key{1, 7},
		'i':  Key{1, 8},
		'o':  Key{1, 9},
		'p':  Key{1, 10},
		'[':  Key{1, 11},
		']':  Key{1, 12},
		'a':  Key{2, 1},
		's':  Key{2, 2},
		'd':  Key{2, 3},
		'f':  Key{2, 4},
		'g':  Key{2, 5},
		'h':  Key{2, 6},
		'j':  Key{2, 7},
		'k':  Key{2, 8},
		'l':  Key{2, 9},
		';':  Key{2, 10},
		'\'': Key{2, 11},
		'z':  Key{3, 1},
		'x':  Key{3, 2},
		'c':  Key{3, 3},
		'v':  Key{3, 4},
		'b':  Key{3, 5},
		'n':  Key{3, 6},
		'm':  Key{3, 7},
		',':  Key{3, 8},
		'.':  Key{3, 9},
		'/':  Key{3, 10},
	}
	nerps = map[rune]Key{
		'x':  Key{1, 1},
		'l':  Key{1, 2},
		'd':  Key{1, 3},
		'p':  Key{1, 4},
		'v':  Key{1, 5},
		'z':  Key{1, 6},
		'k':  Key{1, 7},
		'o':  Key{1, 8},
		'u':  Key{1, 9},
		';':  Key{1, 10},
		'[':  Key{1, 11},
		']':  Key{1, 12},
		'n':  Key{2, 1},
		'r':  Key{2, 2},
		't':  Key{2, 3},
		's':  Key{2, 4},
		'g':  Key{2, 5},
		'y':  Key{2, 6},
		'h':  Key{2, 7},
		'e':  Key{2, 8},
		'i':  Key{2, 9},
		'a':  Key{2, 10},
		'/':  Key{2, 11},
		'j':  Key{3, 1},
		'm':  Key{3, 2},
		'c':  Key{3, 3},
		'w':  Key{3, 4},
		'q':  Key{3, 5},
		'b':  Key{3, 6},
		'f':  Key{3, 7},
		'\'': Key{3, 8},
		',':  Key{3, 9},
		'.':  Key{3, 10},
	}

	L40Charset = map[rune]struct{}{
		// lowercase only, '<' is on the same key as ','; numbers don't move etc.
		'a':  {},
		'b':  {},
		'c':  {},
		'd':  {},
		'e':  {},
		'f':  {},
		'g':  {},
		'h':  {},
		'i':  {},
		'j':  {},
		'k':  {},
		'l':  {},
		'm':  {},
		'n':  {},
		'o':  {},
		'p':  {},
		'q':  {},
		'r':  {},
		's':  {},
		't':  {},
		'u':  {},
		'v':  {},
		'w':  {},
		'x':  {},
		'y':  {},
		'z':  {},
		',':  {},
		'<':  {},
		'.':  {},
		'>':  {},
		'/':  {},
		'?':  {},
		'\'': {},
		'|':  {},
		';':  {},
		':':  {},
		'[':  {},
		'{':  {},
		']':  {},
		'}':  {},
		'-':  {},
		'_':  {},
		'=':  {},
		'+':  {},
	}
)
