package numeron_solver

import "strconv"

type State struct {
	Numbers []rune
}

func (s *State) Check(other *State) (ans Answer) {
	ans.Query = *other
	for i, r := range s.Numbers {
		if r == other.Numbers[i] {
			ans.Hit++
			continue
		}
		for j, p := range other.Numbers {
			if i != j && r == p {
				ans.Blow++
				break
			}
		}
	}
	return
}

func (s *State) IsValid(ans Answer, grade [3]int) bool {
	// https://twitter.com/int0dh/status/876685599490453504
	res := s.Check(&ans.Query)
	sum := grade[0]*2 + grade[1]*3
	for i, n := range s.Numbers {
		if i >= 4 {
			break
		}
		sum += (i + 4) * int(n-'0')
	}
	sum += grade[2] * 8
	sum %= 11
	sum %= 10
	str := strconv.Itoa(sum)
	validMatch := []rune(str)[0] == s.Numbers[4]
	return validMatch && res.Hit == ans.Hit && res.Blow == ans.Blow
}
