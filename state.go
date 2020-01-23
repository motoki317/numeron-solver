package numeron_solver

type State struct {
	Numbers []rune
}

func (s *State) Check(other *State) (ans Answer) {
	ans.Query = *other
	for i, r := range s.Numbers {
		if r == other.Numbers[i] {
			ans.Hit++
		}
		for j, p := range other.Numbers {
			if i != j && r == p {
				ans.Blow++
			}
		}
	}
	return
}

func (s *State) IsValid(ans Answer) bool {
	res := s.Check(&ans.Query)
	return res.Hit == ans.Hit && res.Blow == ans.Blow
}
