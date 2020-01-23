package numeron_solver

type Solver struct {
	Attempts int
	States   []*State
}

func NewSolver(charLength int, charSet []rune) Solver {
	states := make([]*State, 0)

	cases := prepareStates(charLength, charSet, make([]rune, 0), make(map[rune]bool))
	for _, numbers := range cases {
		states = append(states, &State{Numbers:numbers})
	}

	return Solver{States: states}
}

func (s *Solver) RecordAnswer(ans Answer) {
	s.Attempts++
	next := make([]*State, 0)
	for _, state := range s.States {
		if state.IsValid(ans) {
			next = append(next, state)
		}
	}
	s.States = next
}

func prepareStates(remainingLength int, charSet []rune, chars []rune, used map[rune]bool) [][]rune {
	ret := make([][]rune, 0)

	if remainingLength == 0 {
		clone := make([]rune, len(chars))
		copy(clone, chars)
		ret = append(ret, clone)
		return ret
	}

	for _, char := range charSet {
		if isUsed, ok := used[char]; isUsed && ok {
			continue
		}
		used[char] = true
		chars = append(chars, char)
		ret = append(ret, prepareStates(remainingLength - 1, charSet, chars, used)...)
		chars = chars[:len(chars)-1]
		delete(used, char)
	}

	return ret
}
