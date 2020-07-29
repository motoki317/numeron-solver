package simulator

import numeron "github.com/motoki317/numeron-solver"

// When we do not know the answer and want to know all possible remaining combinations
type Simulator interface {
	AskAnswer() (*numeron.Answer, error)
}
