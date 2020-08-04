package numeron_solver

import "strconv"

type Solver struct {
	Attempts        int
	States          []*State
	allowDuplicates bool
	// e.g. 19B -> 1, 9, B(0)
	// B -> 0, M -> 1, D -> 2
	// https://twitter.com/int0dh/status/876685599490453504
	grade [3]int
}

func NewSolver(charLength int, charSet []rune, allowDuplicates bool, gradeString string) Solver {
	states := make([]*State, 0)

	var grade [3]int
	digits, err := strconv.Atoi(gradeString[:2])
	if err != nil {
		panic(err)
	}
	grade[0] = digits / 10
	grade[1] = digits % 10
	switch gradeString[2] {
	case 'B':
		grade[2] = 0
	case 'M':
		grade[2] = 1
	case 'D':
		grade[2] = 2
	default:
		panic("invalid grade")
	}

	cases := prepareStates(charLength, charSet, make([]rune, 0), make(map[rune]bool), allowDuplicates)
	for _, numbers := range cases {
		states = append(states, &State{Numbers: numbers})
	}

	return Solver{States: states, grade: grade}
}

func (s *Solver) RecordAnswer(ans Answer) {
	s.Attempts++
	next := make([]*State, 0)
	for _, state := range s.States {
		if state.IsValid(ans, s.grade) {
			next = append(next, state)
		}
	}
	s.States = next
}

func prepareStates(remainingLength int, charSet []rune, chars []rune, used map[rune]bool, allowDuplicates bool) [][]rune {
	ret := make([][]rune, 0)

	if remainingLength == 0 {
		clone := make([]rune, len(chars))
		copy(clone, chars)
		ret = append(ret, clone)
		return ret
	}

	for _, char := range charSet {
		if isUsed, ok := used[char]; isUsed && ok && !allowDuplicates {
			continue
		}
		chars = append(chars, char)
		ret = append(ret, prepareStates(remainingLength-1, charSet, chars, used, allowDuplicates)...)
		chars = chars[:len(chars)-1]
		delete(used, char)
	}

	return ret
}
