package numeron_solver

type Validator struct {
	Answer State
}

func (v *Validator) GetAnswer(query State) Answer {
	return v.Answer.Check(&query)
}
